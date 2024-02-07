package web

import (
	"go-platform/internal/web/auth"
	"go-platform/internal/web/home"
	"go-platform/pkg/router"
	"go-platform/pkg/server"
)

type Web struct {
	authHandler *auth.AuthHandler
	homeHandler *home.HomeHandler
}

func New() *Web {
	authHandler := auth.New()
	homeHandler := home.New()
	return &Web{
		authHandler,
		homeHandler,
	}

}

func (w *Web) AttachHandler(s *server.Server) {
	s.AttachHandler(func(r *router.Router) {
		w.authHandler.Handle(r)
		w.homeHandler.Handle(r)
	})
}
