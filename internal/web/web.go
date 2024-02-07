package web

import (
	authHandler "go-platform/internal/web/auth"
	baseHandler "go-platform/internal/web/base"
	"go-platform/pkg/router"
	"go-platform/pkg/server"
	"go-platform/pkg/template"
)

type Web struct {
	authHandler *authHandler.AuthHandler
	baseHandler *baseHandler.BaseHandler
}

func New() *Web {
	prefixPath := "web/template"
	template := template.New(prefixPath)
	authHandler := authHandler.New(template)
	baseHandler := baseHandler.New(template)
	return &Web{
		authHandler,
		baseHandler,
	}

}

func (w *Web) AttachHandler(s *server.Server) {
	s.AttachHandler(func(r *router.Router) {
		w.authHandler.Handle(r)
		w.baseHandler.Handle(r)
	})
}
