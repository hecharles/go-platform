package base

import (
	"go-platform/pkg/router"
	"go-platform/pkg/template"
	"net/http"
)

type BaseHandler struct {
	template *template.Template
}

func New(template *template.Template) *BaseHandler {
	return &BaseHandler{template}
}

func (h *BaseHandler) Handle(router *router.Router) {
	prefixPath := ""
	router.Route(prefixPath, h.route)
}

func (h *BaseHandler) route(r *router.Router) {
	r.Get("/", h.index)
}

func (h *BaseHandler) index(w http.ResponseWriter, r *http.Request) error {

	data := struct {
		Title   string
		Content string
	}{
		Title:   "Hello, World!",
		Content: "Welcome to the Go Platform!",
	}

	h.template.Render("index.html", w, data)

	return nil
}
