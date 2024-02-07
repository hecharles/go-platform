package base

import (
	"go-platform/pkg/router"
	"net/http"
)

type BaseHandler struct {
}

func New() *BaseHandler {
	return &BaseHandler{}
}

func (h *BaseHandler) Handle(router *router.Router) {
	prefixPath := ""
	router.Route(prefixPath, h.route)
}

func (h *BaseHandler) route(r *router.Router) {
	r.Get("/", h.index)
}

func (h *BaseHandler) index(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Hello, World!"))

	return nil
}
