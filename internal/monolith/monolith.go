package monolith

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/akaraon8bit/GoEXApp/internal/config"
	"github.com/akaraon8bit/GoEXApp/internal/waiter"
)

type Monolith interface {
	Config() config.AppConfig
	DB() *sql.DB
	Logger() zerolog.Logger
	Mux() *chi.Mux
	RPC() *grpc.Server // Can return nil for non-gRPC deployments
	Waiter() waiter.Waiter
}

type Module interface {
	Startup(context.Context, Monolith) error
}
