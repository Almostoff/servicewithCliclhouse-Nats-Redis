package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"taskFive/server/internal/cases"
	"taskFive/server/internal/entity"
)

type Controller struct {
	usecase cases.Usecase
}

func NewController(usecase cases.Usecase) *Controller {
	return &Controller{
		usecase: usecase,
	}
}

func Build(r *chi.Mux, usecase cases.Usecase) {
	ctr := NewController(usecase)

	r.Get("/items/list", ctr.GetAll)

	r.Route("/{campaignID}", func(r chi.Router) {
		r.Post("/item/create", ctr.CreateItem)
		r.Patch("/item/update/{id}", ctr.UpdateItem)
		r.Delete("/item/remove/{id}", ctr.DeleteItem)
	})

}

func (s *Controller) CreateItem(w http.ResponseWriter, r *http.Request) {
	campId := chi.URLParam(r, "campaignID")
	if campId == "" {
		log.Fatal("invalid campaign id: ")
		return
	}
	item := &entity.Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	cId, _ := strconv.Atoi(campId)
	data, err := s.usecase.CreateItem(item, cId)
	a, err := json.Marshal(data)
	w.Write(a)
}

func (s *Controller) UpdateItem(w http.ResponseWriter, r *http.Request) {
	campId := chi.URLParam(r, "campaignID")
	iID := chi.URLParam(r, "id")

	item := &entity.Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	cId, _ := strconv.Atoi(campId)
	iId, _ := strconv.Atoi(iID)

	data, err := s.usecase.PatchItem(item, cId, iId)
	if data == nil {
		w.WriteHeader(404)
		w.Write([]byte("errors.item.notFound\n"))
	}
	if err != nil {
		log.Println("can`t update item: ", err)
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Controller) DeleteItem(w http.ResponseWriter, r *http.Request) {
	campId := chi.URLParam(r, "campaignID")
	iID := chi.URLParam(r, "id")

	if campId == "" || iID == "" {
		log.Fatal("invalid campaign or item id: ")
		return
	}

	cId, _ := strconv.Atoi(campId)
	iId, _ := strconv.Atoi(iID)

	data, err := s.usecase.DeleteItem(iId, cId)
	if data == nil {
		w.WriteHeader(404)
		w.Write([]byte("errors.item.notFound\n"))
	}
	if err != nil {
		log.Println("can`t delete item: ", err)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := s.usecase.GetAll()
	if err != nil {
		log.Fatal("can`t get items list", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
