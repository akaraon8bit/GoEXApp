package application

import (
	"context"

	"github.com/akaraon8bit/GoEXApp/notifications/internal/models"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}
