package home

import (
	"go-platform/pkg/router"
	"net/http"
)

type HomeHandler struct {
}

func New() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Handle(router *router.Router) {
	prefixPath := ""
	router.Route(prefixPath, h.route)
}

func (h *HomeHandler) route(r *router.Router) {
	r.Get("/", h.home)
}

func (h *HomeHandler) home(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Hello, World!"))

	return nil
}
