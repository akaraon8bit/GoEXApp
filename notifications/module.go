package notifications

import (
	"context"
	// "net/http"

	"github.com/akaraon8bit/GoEXApp/internal/monolith"
	"github.com/akaraon8bit/GoEXApp/notifications/internal/application"
	"github.com/akaraon8bit/GoEXApp/notifications/internal/logging"
	"github.com/akaraon8bit/GoEXApp/notifications/internal/rest"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	// setup application
	var app application.App
	app = application.New(nil) // For notifications, we might not need customer repo
	app = logging.LogApplicationAccess(app, mono.Logger())

	// setup REST handlers
	handlers := rest.NewNotificationsHandlers(app)
	handlers.RegisterRoutes(mono.Mux())

	return nil
}
