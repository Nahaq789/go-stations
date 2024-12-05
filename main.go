package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		<-ctx.Done()

		fmt.Println("シグナル受信")

		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(c); err != nil {
			log.Printf("Server shutdown error: %v", err)
			return
		}
	}(ctx)
	log.Println(server.ListenAndServe())
	wg.Wait()

	return nil

	// チャネル作成
	// シャットダウン時にgoroutineの関数より先に
	// main関数が終了するの防ぐ
	// idleConnsClosed := make(chan struct{})
	// go func() {
	// 	// os.Signal型のチャネル作成
	// 	sigint := make(chan os.Signal, 1)
	// 	// osからのSIGINTとSIGTERMを受け取れるように設定
	// 	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	// 	// チャネルから値を受信
	// 	// 値を受信するまで処理ストップ
	// 	<-sigint

	// 	// シャットダウンしてから10秒間は継続する
	// 	c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// 	defer cancel()

	// 	// チャネルから値を受信したら
	// 	// サーバーシャットダウン
	// 	if err := server.Shutdown(c); err != nil {
	// 		log.Printf("Server shutdown error: %v", err)
	// 		close(idleConnsClosed)
	// 		return
	// 	}

	// 	log.Println("Server is shut down")
	// 	// 正常終了時もチャネルを閉じる
	// 	close(idleConnsClosed)
	// }()

	// log.Printf("Server is running on %s", server.Addr)
	// if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
	// 	log.Fatalf("HTTP server error: %v", err)
	// }

	// <-idleConnsClosed
}
