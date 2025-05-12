package rest

import (
	// "context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/akaraon8bit/GoEXApp/notifications/notificationspb"
	"github.com/akaraon8bit/GoEXApp/notifications/internal/application"
)

type NotificationsHandlers struct {
	app application.App
}

func NewNotificationsHandlers(app application.App) *NotificationsHandlers {
	return &NotificationsHandlers{app: app}
}

func (h *NotificationsHandlers) RegisterRoutes(r *chi.Mux) {
	r.Post("/api/notifications/order-created", h.notifyOrderCreated)
	r.Post("/api/notifications/order-canceled", h.notifyOrderCanceled)
	r.Post("/api/notifications/order-ready", h.notifyOrderReady)
}

func (h *NotificationsHandlers) notifyOrderCreated(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req notificationspb.NotifyOrderCreatedRequest
	if err := decodeRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.app.NotifyOrderCreated(ctx, application.OrderCreated{
		OrderID:    req.GetOrderId(),
		CustomerID: req.GetCustomerId(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeResponse(w, r, http.StatusOK, &notificationspb.NotifyOrderCreatedResponse{})
}

func (h *NotificationsHandlers) notifyOrderCanceled(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req notificationspb.NotifyOrderCanceledRequest
	if err := decodeRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.app.NotifyOrderCanceled(ctx, application.OrderCanceled{
		OrderID:    req.GetOrderId(),
		CustomerID: req.GetCustomerId(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeResponse(w, r, http.StatusOK, &notificationspb.NotifyOrderCanceledResponse{})
}

func (h *NotificationsHandlers) notifyOrderReady(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req notificationspb.NotifyOrderReadyRequest
	if err := decodeRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.app.NotifyOrderReady(ctx, application.OrderReady{
		OrderID:    req.GetOrderId(),
		CustomerID: req.GetCustomerId(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeResponse(w, r, http.StatusOK, &notificationspb.NotifyOrderReadyResponse{})
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
