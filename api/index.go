package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"golang.org/x/sync/errgroup"

	"github.com/akaraon8bit/GoEXApp/customers"
	"github.com/akaraon8bit/GoEXApp/internal/config"
	"github.com/akaraon8bit/GoEXApp/internal/logger"
	"github.com/akaraon8bit/GoEXApp/internal/monolith"
	"github.com/akaraon8bit/GoEXApp/internal/waiter"
	"github.com/akaraon8bit/GoEXApp/notifications"
	"database/sql"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type vercelMonolith struct {
	cfg    config.AppConfig
	db     *sql.DB
	logger zerolog.Logger
	mux    *chi.Mux
	waiter waiter.Waiter
	rpc    *grpc.Server // Added but not initialized for Vercel
}

func (v *vercelMonolith) Config() config.AppConfig {
	return v.cfg
}

func (v *vercelMonolith) DB() *sql.DB {
	return v.db
}

func (v *vercelMonolith) Logger() zerolog.Logger {
	return v.logger
}

func (v *vercelMonolith) Mux() *chi.Mux {
	return v.mux
}

func (v *vercelMonolith) RPC() *grpc.Server {
	// Return nil for Vercel deployment
	return nil
}

func (v *vercelMonolith) Waiter() waiter.Waiter {
	return v.waiter
}

func initMonolith() (*vercelMonolith, error) {
	var cfg config.AppConfig
	var err error

	// Initialize configuration
	cfg, err = config.InitConfig()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	logger := logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    logger.Level(cfg.LogLevel),
	})

	// Initialize database if needed
	var db *sql.DB
	if cfg.PG.Conn != "" {
		db, err = sql.Open("pgx", cfg.PG.Conn)
		if err != nil {
			return nil, err
		}
	}

	// Create router with CORS
	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Create waiter
	waiter := waiter.New(waiter.CatchSignals())

	return &vercelMonolith{
		cfg:    cfg,
		db:     db,
		logger: logger,
		mux:    mux,
		waiter: waiter,
		rpc:    nil, // Explicitly set to nil for Vercel
	}, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize monolith
	mono, err := initMonolith()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Initialize modules
	modules := []monolith.Module{
		&customers.Module{},
		&notifications.Module{},
	}

	// Start modules
	g, ctx := errgroup.WithContext(r.Context())
	for _, module := range modules {
		module := module
		g.Go(func() error {
			return module.Startup(ctx, mono)
		})
	}

	if err := g.Wait(); err != nil {
		http.Error(w, "Module initialization failed", http.StatusInternalServerError)
		return
	}

	// Serve the request
	mono.Mux().ServeHTTP(w, r)
}
