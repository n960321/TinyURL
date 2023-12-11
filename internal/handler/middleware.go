package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		log.Info().
			Str("URI", req.RequestURI).
			Any("Header", req.Header).
			Any("Body", string(body)).
			Msgf("Request comming")

		next.ServeHTTP(resp, req)
	})
}
