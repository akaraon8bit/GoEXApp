package rest

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/akaraon8bit/GoEXApp/customers/customerspb"
	"github.com/akaraon8bit/GoEXApp/customers/internal/application"
	"github.com/google/uuid"
)

type CustomersHandlers struct {
	app application.App
}

func NewCustomersHandlers(app application.App) *CustomersHandlers {
	return &CustomersHandlers{app: app}
}

func (h *CustomersHandlers) RegisterRoutes(r *chi.Mux) {
	r.Post("/api/customers", h.registerCustomer)
	r.Get("/api/customers/{id}", h.getCustomer)
	r.Put("/api/customers/{id}/enable", h.enableCustomer)
	r.Put("/api/customers/{id}/disable", h.disableCustomer)
	r.Post("/api/customers/{id}/authorize", h.authorizeCustomer)
}

func (h *CustomersHandlers) registerCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req customerspb.RegisterCustomerRequest
	if err := decodeRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	err := h.app.RegisterCustomer(ctx, application.RegisterCustomer{
		ID:        id,
		Name:      req.GetName(),
		SmsNumber: req.GetSmsNumber(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &customerspb.RegisterCustomerResponse{Id: id}
	encodeResponse(w, r, http.StatusCreated, res)
}

func (h *CustomersHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	customer, err := h.app.GetCustomer(ctx, application.GetCustomer{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := &customerspb.GetCustomerResponse{
		Customer: &customerspb.Customer{
			Id:        customer.ID(),
			Name:      customer.Name,
			SmsNumber: customer.SmsNumber,
			Enabled:   customer.Enabled,
		},
	}
	encodeResponse(w, r, http.StatusOK, res)
}

func (h *CustomersHandlers) enableCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.app.EnableCustomer(ctx, application.EnableCustomer{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeResponse(w, r, http.StatusOK, &customerspb.EnableCustomerResponse{})
}

func (h *CustomersHandlers) disableCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.app.DisableCustomer(ctx, application.DisableCustomer{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeResponse(w, r, http.StatusOK, &customerspb.DisableCustomerResponse{})
}

func (h *CustomersHandlers) authorizeCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.app.AuthorizeCustomer(ctx, application.AuthorizeCustomer{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeResponse(w, r, http.StatusOK, &customerspb.AuthorizeCustomerResponse{})
}

func decodeRequest(r *http.Request, msg proto.Message) error {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return protojson.Unmarshal(body, msg)
}

func encodeResponse(w http.ResponseWriter, r *http.Request, status int, msg proto.Message) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	m := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}

	jsonBytes, err := m.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonBytes)
}
