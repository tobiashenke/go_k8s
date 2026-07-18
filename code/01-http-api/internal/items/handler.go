package items

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ItemHandler struct {
	service   *ItemService
	validator *validator.Validate
}

func NewItemHandler(i *ItemService) *ItemHandler {
	return &ItemHandler{
		service:   i,
		validator: validator.New(),
	}
}

func (h *ItemHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAll(r.Context())
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
	err = h.validator.Struct(item)
	if err != nil {
		writeError(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}
	err = h.service.Save(r.Context(), item)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save item")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "item not deleted")
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
