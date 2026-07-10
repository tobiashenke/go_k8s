package items

import (
	"encoding/json"
	"net/http"
)

type ItemHandler struct {
	ItemService *ItemService
}

func NewItemHandler(i *ItemService) *ItemHandler {
	return &ItemHandler{ItemService: i}
}

func (h *ItemHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.ItemService.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve items")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = h.ItemService.Save(item)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save item")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.ItemService.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

type ErrorResponse struct {
	Type   string
	Title  string
	Status int
	Detail string
}

func writeError(w http.ResponseWriter, sc int, mes string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(sc)
	e := &ErrorResponse{
		Type:   "about:blank",
		Title:  http.StatusText(sc),
		Status: sc,
		Detail: mes,
	}
	_ = json.NewEncoder(w).Encode(e)
}
