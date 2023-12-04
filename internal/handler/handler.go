package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tinyurl/pkg/base58"

	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{url_key}", h.redirection).Methods(http.MethodGet)

	sr := r.PathPrefix("/api/v1").Subrouter()
	sr.HandleFunc("/create", h.create).Methods(http.MethodPost)
	sr.HandleFunc("/redirection/{url_key}", h.redirectionWithHttpResp).Methods(http.MethodGet)
	return r
}



type CreateReq struct {
	Url string `json:"url"`
}


func (h *Handler) create(rw http.ResponseWriter, req *http.Request) {
	var createReq CreateReq

	err := json.NewDecoder(req.Body).Decode(&createReq)
	if err != nil {
		log.Printf("unmarshall req failed, err: %v", err)
		return
	}

	fmt.Printf("createReq: %+v\n", createReq)
}

func (h *Handler) redirection(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	urlKey := v["url_key"]

	decodeValue, err := base58.DecodeToInt(urlKey)
	if err != nil {
		log.Printf("the url_key is invaild, err:%v \n", err)
		return
	}

	fmt.Printf("decodeValue: %v\n", decodeValue)
	http.Redirect(rw, req, "https://www.google.com", http.StatusSeeOther)
}

func (h *Handler) redirectionWithHttpResp(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	urlKey := v["url_key"]

	decodeValue, err := base58.DecodeToInt(urlKey)
	if err != nil {
		log.Printf("the url_key is invaild, err:%v \n", err)
		return
	}
	fmt.Printf("decodeValue: %v\n", decodeValue)
}
