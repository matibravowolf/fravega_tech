package domain

import (
	"time"

	"github.com/google/uuid"
)

type RouteStatus string

const (
	RouteStatusPending   RouteStatus = "pending"   // Route created, not started yet
	RouteStatusOnWay     RouteStatus = "on_way"    // Driver is on the way
	RouteStatusCompleted RouteStatus = "completed" // All deliveries completed
	RouteStatusCancelled RouteStatus = "cancelled" // Route was cancelled
)

type Route struct {
	ID        string      `json:"id"`
	Vehicle   string      `json:"vehicle"`
	Driver    string      `json:"driver"`
	Status    RouteStatus `json:"status"`
	Purchases []*Purchase `json:"purchases"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewRoute(vehicle, driver string) *Route {
	return &Route{
		ID:        uuid.New().String(),
		Vehicle:   vehicle,
		Driver:    driver,
		Purchases: make([]*Purchase, 0),
		Status:    RouteStatusPending,
		CreatedAt: time.Now(),
	}
}

func (r *Route) AddPurchase(purchase *Purchase) error {
	if err := purchase.ValidatePurchaseAssignment(purchase, r); err != nil {
		return err
	}

	r.Purchases = append(r.Purchases, purchase)

	return nil
}

func (r *Route) IsPendingRoute() bool {
	return r.Status == RouteStatusPending
}

func (r *Route) MarkPurchaseAsDelivered(purchaseID string) error {
	purchase, found := r.findPurchaseByID(purchaseID)
	if !found {
		return NewError(ErrorPurchaseNotExist, "purchase does not exist in route", nil).
			WithMetadata(RouteID, r.ID)
	}

	purchase.Status = PurchaseDelivered
	return nil
}

func (r *Route) findPurchaseByID(purchaseID string) (*Purchase, bool) {
	for i := range r.Purchases {
		if r.Purchases[i].ID == purchaseID {
			return r.Purchases[i], true
		}
	}
	return nil, false
}
