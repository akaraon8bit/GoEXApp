package customers

import (
	"context"

	"github.com/akaraon8bit/GoEXApp/customers/internal/application"
	"github.com/akaraon8bit/GoEXApp/customers/internal/grpc"
	"github.com/akaraon8bit/GoEXApp/customers/internal/logging"
	"github.com/akaraon8bit/GoEXApp/customers/internal/postgres"
	"github.com/akaraon8bit/GoEXApp/customers/internal/rest"
	"github.com/akaraon8bit/GoEXApp/internal/ddd"
	"github.com/akaraon8bit/GoEXApp/internal/monolith"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
	customers := postgres.NewCustomerRepository("customers.customers", mono.DB())

	// setup application
	app := logging.LogApplicationAccess(
		application.New(customers, domainDispatcher),
		mono.Logger(),
	)

	// Register gRPC server if RPC is available, otherwise register REST handlers directly
	if rpc := mono.RPC(); rpc != nil {
		// gRPC mode
		if err := grpc.RegisterServer(app, rpc); err != nil {
			return err
		}
		// Register REST gateway
		if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
			return err
		}
	} else {
		// HTTP-only mode (for Vercel)
		handlers := rest.NewCustomersHandlers(app)
		handlers.RegisterRoutes(mono.Mux())
	}

	// Register Swagger UI in both modes
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	return nil
}
