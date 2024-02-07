package router

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type (
	// HandlerFunc is a modified version of http.HandlerFunc with error returns.
	// The standard http.HandlerFunc doesn't return any error, and errors from
	// the handler cannot be propogated easily to the middleware(for logging, etc).
	HandlerFunc func(w http.ResponseWriter, r *http.Request) error
	// MiddlewareFunc is a middleware function that compatible with HandlerFunc.
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	Handler interface {
		PrefixPath() string
		Handle(r *Router)
	}
)

// ToHTTPHandler convert the router.HandlerFunc to standard http.HandlerFunc.
// Please note the error return will be ignored when its converted to standard
// http.HandlerFunc.
func (h HandlerFunc) ToHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}
}

// ChainedMiddleware is a stack of middleware to used as an addition to middlewares
// in the router.
type ChainedMiddleware struct {
	middlewares []MiddlewareFunc
	r           *Router
}

// ChainMiddlewares returns a chained middleware that stacked and executed in FIFO order.
//
// For example: middleware_1(middleware_2(handler)).
// The middleware_2 will wrap the handler first then middleware_1 will wrap middleware_2 and handler.
func ChainMiddlewares(r *Router, middlewares ...MiddlewareFunc) *ChainedMiddleware {
	return &ChainedMiddleware{
		middlewares: middlewares,
		r:           r,
	}
}

// Then for serving request based on the Method requriment.
func (c *ChainedMiddleware) Then(method, path string, handler HandlerFunc) {
	for i := range c.middlewares {
		handler = c.middlewares[len(c.middlewares)-i-1](handler)
	}
	c.r.handleRoute(method, path, handler)
}

// Get for serving 'GET' request based on the middleware chain.
func (c *ChainedMiddleware) Get(path string, handler HandlerFunc) {
	c.Then(http.MethodGet, path, handler)
}

// Patch for serving 'PATCH' request based on the middleware chain.
func (c *ChainedMiddleware) Patch(path string, handler HandlerFunc) {
	c.Then(http.MethodPatch, path, handler)
}

// Post for serving 'POST' request based on the middleware chain.
func (c *ChainedMiddleware) Post(path string, handler HandlerFunc) {
	c.Then(http.MethodPost, path, handler)
}

// Put for serving 'PUT' request based on the middleware chain.
func (c *ChainedMiddleware) Put(path string, handler HandlerFunc) {
	c.Then(http.MethodPut, path, handler)
}

// Delete for serving 'DELETE' request based on the middleware chain.
func (c *ChainedMiddleware) Delete(path string, handler HandlerFunc) {
	c.Then(http.MethodDelete, path, handler)
}

// Head for serving 'HEAD' request based on the middleware chain.
func (c *ChainedMiddleware) Head(path string, handler HandlerFunc) {
	c.Then(http.MethodHead, path, handler)
}

// Options for serving 'OPTIONS' request based on the middleware chain.
func (c *ChainedMiddleware) Options(path string, handler HandlerFunc) {
	c.Then(http.MethodOptions, path, handler)
}

// Handler returns the standard http.Handler using the ChainedMiddleware
// and HandlerFunc function.
func (c *ChainedMiddleware) Handler(handler HandlerFunc) http.Handler {
	for i := range c.middlewares {
		handler = c.middlewares[len(c.middlewares)-i-1](handler)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
}

type Mux interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	Handle(path string, handler http.Handler) *mux.Route
	Walk(walkfn mux.WalkFunc) error
}

// Router is a wrapper around httprouter to introduce middleware chaining and
// error propagation from the http handler.
type Router struct {
	prefixPath  string
	r           Mux
	middlewares []MiddlewareFunc
}

// New creates a new router.
func New() *Router {
	router := mux.NewRouter()
	return &Router{
		r: router,
	}
}

// Route makes it possible to create a route group.
func (r *Router) Route(prefix string, fn func(r *Router)) {

	if r.prefixPath != "" {
		prefix = r.prefixPath + prefix
	}
	copyR := &Router{
		r:          r.r,
		prefixPath: prefix,
	}
	fn(copyR)
}

func (r *Router) PrefixPath() string {
	return r.prefixPath
}

// Use will appends a global middleware to the router. The appended middleware are executed based on the registration stack.
// The first registered middleware will be executed first or in other words, in FIFO manner.
//
// For example: middleware_1(middleware_2(handler)).
// The middleware_2 will wrap the handler first then middleware_1 will wrap middleware_2 and handler.
func (r *Router) Use(middlewares ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// Get for serving 'GET' request.
func (r *Router) Get(path string, handler HandlerFunc) {
	r.handleRoute(http.MethodGet, path, handler)
}

// Patch for serving 'PATCH' request.
func (r *Router) Patch(path string, handler HandlerFunc) {
	r.handleRoute(http.MethodPatch, path, handler)
}

// Post for serving 'POST' request.
func (r *Router) Post(path string, handler HandlerFunc) {
	r.handleRoute(http.MethodPost, path, handler)
}

// Put for serving 'PUT' request.
func (r *Router) Put(path string, handler HandlerFunc) {
	r.handleRoute(http.MethodPut, path, handler)
}

// Delete for serving 'DELETE' request.
func (r *Router) Delete(path string, handler HandlerFunc) {
	r.handleRoute(http.MethodDelete, path, handler)
}

// Head for serving 'HEAD' request.
func (r *Router) Head(path string, handler HandlerFunc) {
	r.handleRoute(http.MethodHead, path, handler)
}

// Options for serving 'OPTIONS' request.
func (r *Router) Options(path string, handler HandlerFunc) {
	r.handleRoute(http.MethodOptions, path, handler)
}

func (r *Router) Handle(path string, handler http.Handler) {
	r.handle(path, handler)
}

// handle will handle all http request for all http method. Handle is different from handlerFunc
// because it does not care about what method is assigned to a specific path.
//
// nolint: all
func (r *Router) handle(p string, handler http.Handler) {
	if r.prefixPath != "" {
		path.Join(r.prefixPath, p)
	}

	h := func(w http.ResponseWriter, r *http.Request) error {
		handler.ServeHTTP(w, r)
		return nil
	}
	for i := range r.middlewares {
		h = r.middlewares[len(r.middlewares)-1-i](h)
	}

	handlerFunc := func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			r.handlePanic(req.Context(), w, p, err)
		}
	}

	r.r.Handle(p, http.HandlerFunc(handlerFunc))
}

func (r *Router) handleRoute(method, p string, handler HandlerFunc) {
	if r.prefixPath != "" {
		p = r.prefixPath + p
	}
	h := handler

	// Stack the middleware from the last one to the frist one. But because we are stacking/wrapping them backwards,
	// the first middleware will be the first one to be executed as the (N) middleware will be wrapped with (N-1).
	for i := range r.middlewares {
		h = r.middlewares[len(r.middlewares)-1-i](h)
	}

	handlerFunc := func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			r.handlePanic(req.Context(), w, p, err)
		}

	}

	// OPTIONS is added to handle preflight CORS request from browsers.
	if method != http.MethodOptions {
		r.r.HandleFunc(p, handlerFunc).Methods(method, http.MethodOptions)
	} else {
		r.r.HandleFunc(p, handlerFunc).Methods(method)
	}
}

func (r *Router) handlePanic(ctx context.Context, w http.ResponseWriter, path string, v interface{}) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "internal error"}`))

	panicErr, ok := v.(error)
	if !ok {
		panicErr = fmt.Errorf("recover from panic: %v", v)

	} else {
		panicErr = fmt.Errorf("recover from panic: %w", panicErr)

	}

	return panicErr
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

// Vars retrieves route variables/path parameters, if any.
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
