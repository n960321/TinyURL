package handler

import (
	"encoding/json"
	"net/http"
	"tinyurl/internal/service"
	"tinyurl/pkg/base58"
	"tinyurl/pkg/database"
	redispkg "tinyurl/pkg/redis"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc service.URLGenerateServicer
}

func NewHandler(db *database.Database, cache *redispkg.RedisCache) *Handler {
	return &Handler{
		svc: service.NewURLGenerateService(db, cache),
	}
}

func (h *Handler) GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path("/metrics").Handler(promhttp.Handler())
	r.Use(loggingMiddleware)
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
	var resp = new(CreateResp)
	err := json.NewDecoder(req.Body).Decode(&createReq)
	if err != nil {
		log.Warn().Err(err).Msg("create: Unmarshall req failed")
		return
	}

	urlInfo := &service.UrlInfo{URL: createReq.Url}
	urlInfo, err = h.svc.CreateURLInfo(urlInfo)
	if err != nil {
		log.Warn().Err(err).Msg("create: CreateURLInfo failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return

	}
	resp.UrlKey = base58.EncodeFromInt(int(urlInfo.ID))
	rw.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Warn().Err(err).Msg("create: marshal failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(respBody)
}

func (h *Handler) redirection(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	urlKey := v["url_key"]

	id, err := base58.DecodeToInt(urlKey)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	urlInfo, err := h.svc.GetURLInfo(id)
	if err != nil {
		log.Warn().Err(err).Msg("redirection: GetURLInfo failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, req, urlInfo.URL, http.StatusSeeOther)
}

type RedirectionWithHttpResp struct {
	Url string `json:"url"`
}

func (h *Handler) redirectionWithHttpResp(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	urlKey := v["url_key"]
	resp := &RedirectionWithHttpResp{}
	id, err := base58.DecodeToInt(urlKey)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	urlInfo, err := h.svc.GetURLInfo(id)
	if err != nil {
		log.Warn().Err(err).Msg("redirection: GetURLInfo failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.Url = urlInfo.URL
	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Warn().Err(err).Msg("redirection: marshal failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(respBody)
}
