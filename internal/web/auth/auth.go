package auth

import (
	"go-platform/pkg/router"
	"net/http"
)

type AuthHandler struct {
}

func New() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Handle(router *router.Router) {
	prefixPath := "/auth"
	router.Route(prefixPath, h.route)
}

func (h *AuthHandler) route(r *router.Router) {
	r.Get("/login", h.login)
	r.Get("/register", h.register)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Hello, login!"))

	return nil
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("Hello, register!"))

	return nil
}
