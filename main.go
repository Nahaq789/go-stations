package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	// TODO: サーバーをlistenする
	server := http.Server{
		Addr:    port,
		Handler: mux,
		// Handler: middleware.Recovery(mux),
	}

	// チャネル作成
	idleConnsClosed := make(chan struct{})
	go func() {
		// os.Signal型のチャネル作成
		sigint := make(chan os.Signal, 1)
		// osからのSIGINTとSIGTERMを受け取れるように設定
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		// チャネルから値を受信
		// 値を受信するまで処理ストップ
		<-sigint

		// チャネルから値を受信したら
		// サーバーシャットダウン
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("server shutdown error: %v", err)
			return
		}
		log.Printf("Server is shut down")
		close(idleConnsClosed)
	}()

	log.Printf("Server is running on %s", server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}

	<-idleConnsClosed
	return nil
}
