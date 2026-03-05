package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/category"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Server struct {
	host string
	port string

	httpsrv *http.Server
	db      *database.Database

	logger *slog.Logger

	errc chan error
	ctx  context.Context
}

type ServerOpts struct {
	Host     string
	Port     string
	Database *database.Database
	Logger   *slog.Logger
	Context  context.Context
}

// Creates a new server initializing the endpoint structure and abstracts some of the functionality
// behind a set of functions
func New(opts *ServerOpts) *Server {
	// in a different scenario I'd validate the input, it all depends if I'm exposing this interface
	// to someone else as a library or even a simple interface.
	s := &Server{
		host:   opts.Host,
		port:   opts.Port,
		db:     opts.Database,
		logger: opts.Logger,
		ctx:    opts.Context,
		errc:   make(chan error),
	}

	//TODO(rafael.nunes): change the way repositories are initialized and given to handlers

	prodRepo := models.NewProductsRepository(s.db.GORM())
	cateRepo := models.NewCategoryRepository(s.db.GORM())

	cata := catalog.NewCatalogHandler(&catalog.CatalogHandlerOpts{Repository: prodRepo, Logger: opts.Logger})
	cate := category.NewCatalogHandler(&category.CategoryHandlerOpts{Repository: cateRepo, Logger: opts.Logger})

	// set up routing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", cata.HandleGet)
	mux.HandleFunc("GET /catalog/{code}", cata.HandleGetByCode)

	mux.HandleFunc("GET /categories", cate.HandleGet)
	mux.HandleFunc("POST /categories", cate.HandleCreate)

	// set up the HTTP server
	s.httpsrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.host, s.port),
		Handler: mux,
	}

	return s
}

func (s *Server) Start() {
	s.logger.Info("starting server", slog.String("address", s.host+":"+s.port))

	go func() {
		if err := s.httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errc <- err
		}
	}()
}

func (s *Server) Stop() {
	if err := s.httpsrv.Shutdown(s.ctx); err != nil {
		s.logger.Error("failed to shutdown server gracefully, will be stopped either way", slog.Any("error", err))
	}
}

func (s *Server) Error() <-chan error {
	return s.errc
}
