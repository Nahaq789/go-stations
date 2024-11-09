package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (t *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var m *model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}

		if m.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var result *model.CreateTODOResponse
		result, err = t.Create(r.Context(), m)
		if err != nil {
			log.Println(err)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodPut:
		var m *model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}

		if m.ID == 0 || m.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := t.Update(r.Context(), m)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodGet:
		prevId, _ := strconv.Atoi(r.FormValue("prev_id"))
		size, _ := strconv.Atoi(r.FormValue("size"))
		m := model.ReadTODORequest{
			PrevID: int64(prevId),
			Size:   int64(size),
		}
		result, err := t.Read(r.Context(), &m)
		if err != nil {
			log.Println(err)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodDelete:
		var m *model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}
		if len(m.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := t.Delete(r.Context(), m)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, _ := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	res := model.CreateTODOResponse{
		TODO: *todo,
	}

	return &res, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	t := []model.TODO{}
	todos, _ := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	for _, todo := range todos {
		t = append(t, *todo)
	}

	res := model.ReadTODOResponse{
		TODOs: t,
	}
	return &res, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, _ := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	if todo == nil {
		return nil, &model.ErrNotFound{
			When: time.Now(),
			What: "hogehoge",
		}
	}
	res := model.UpdateTODOResponse{
		TODO: *todo,
	}
	return &res, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, &model.ErrNotFound{
			When: time.Now(),
			What: "hogehoge",
		}
	}
	res := model.DeleteTODOResponse{}
	return &res, nil
}
