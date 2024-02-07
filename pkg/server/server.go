package server

import (
	"context"
	"errors"
	"fmt"
	"go-platform/pkg/router"
	log "log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	defaultServerAddress = ":8080"

	defaultIdleTimeout  = time.Second * 5
	defaultReadTimeout  = time.Second * 10
	defaultWriteTimeout = time.Second * 5
)

// AppInfo is an interface to define '/version' endpoint for http server
// which export the application identifier.
type AppInfo interface {
	Name() string
	Version() string
	Commit() string
}

type Config struct {
	Address      string        `yaml:"address"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

func (c *Config) Validate() error {
	if c.Address == "" {
		c.Address = defaultServerAddress
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = time.Duration(defaultIdleTimeout)
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = time.Duration(defaultReadTimeout)
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = time.Duration(defaultWriteTimeout)
	}
	return nil
}

type Server struct {
	mainServer   *http.Server
	mainListener net.Listener
	// gRPCHandler is used for handling grpc-gateway requests.

	httpHandler     *router.Router
	handlerAttacher []func(r *router.Router)

	mu      sync.Mutex
	running bool
}

func New(c Config) (*Server, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	mainListener, err := createListener("tcp", c.Address)
	if err != nil {
		return nil, err
	}

	s := &Server{

		mainServer: &http.Server{
			IdleTimeout:       time.Duration(c.IdleTimeout),
			ReadHeaderTimeout: time.Duration(c.ReadTimeout),
			WriteTimeout:      time.Duration(c.WriteTimeout),
		},
		mainListener: mainListener,
		httpHandler:  router.New(),
	}

	return s, nil
}

func createListener(network, addr string) (net.Listener, error) {
	var (
		listener net.Listener
		err      error
	)

	// if the listener is not iniherited, create a new listener
	// from net listener instead of using the upgrader listener.
	if listener == nil {
		listener, err = net.Listen(network, addr)
	}

	return listener, err
}

func (s *Server) Name() string {
	return "http-server"
}

func (s *Server) Listener() net.Listener {
	return s.mainListener
}

func (s *Server) AttachHandler(fn func(r *router.Router)) {
	s.handlerAttacher = append(s.handlerAttacher, fn)
}

func (s *Server) serve(ctx context.Context) error {
	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	// Attach all handler before serving requests.
	for _, attacher := range s.handlerAttacher {
		attacher(s.httpHandler)
	}

	s.mainServer.Handler = s.httpHandler

	s.mu.Lock()
	s.running = true
	s.mu.Unlock()

	if err := s.mainServer.Serve(s.mainListener); err != nil {
		return fmt.Errorf("http-server: error when serving: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info("shutting down http server...")
	return s.mainServer.Shutdown(ctx)
}

func (s *Server) ready(ctx context.Context) error {
	s.mu.Lock()
	ready := s.running
	s.mu.Unlock()

	if ready {
		return nil
	}
	return errors.New("http server is not running")
}

func (s *Server) Start(ctx context.Context) error {

	go func(ctx context.Context) {
		if err := s.serve(ctx); err != nil {
			log.Error("http server: %v", err)
		}
	}(ctx)
	for {
		if err := s.ready(ctx); err == nil {
			log.Info("Http server is ready", "running on", s.mainListener.Addr().String())
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}
