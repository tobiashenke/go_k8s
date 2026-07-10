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
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ItemService.Save(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}
