package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"tinyurl/internal/service"
	"tinyurl/pkg/database"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc service.URLGenerateServicer
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{
		svc: service.NewURLGenerateService(db),
	}
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

type CreateResp struct {
	UrlKey string `json:"url_key"`
}

func (h *Handler) create(rw http.ResponseWriter, req *http.Request) {
	var createReq CreateReq

	err := json.NewDecoder(req.Body).Decode(&createReq)
	if err != nil {
		log.Printf("unmarshall req failed, err: %v", err)
		return
	}
	log.Printf("createReq: %+v\n", createReq)

	urlKey := h.svc.CreateShortURL(createReq.Url)
	resp := CreateResp{UrlKey: urlKey}
	rw.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Printf("marshal failed, err: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(respBody)
}

func (h *Handler) redirection(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	urlKey := v["url_key"]

	// TODO: check existing in cache
	url := h.svc.GetShortURL(urlKey)
	http.Redirect(rw, req, url, http.StatusSeeOther)
}

type RedirectionWithHttpResp struct {
	Url string `json:"url"`
}

func (h *Handler) redirectionWithHttpResp(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	urlKey := v["url_key"]

	// TODO: check existing in cache
	url := h.svc.GetShortURL(urlKey)

	resp := RedirectionWithHttpResp{
		Url: url,
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Printf("marshal failed, err:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(respBody)
}
