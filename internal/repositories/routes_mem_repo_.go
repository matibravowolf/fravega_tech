package repositories

import (
	"context"
	"sync"

	"github.com/uMakeMeCrazy/fravega_tech/internal/core/domain"
)

type MemoryRepo struct {
	mu     sync.RWMutex
	routes map[string]*domain.Route
}

func NewMemoryRepository() *MemoryRepo {
	return &MemoryRepo{
		routes: make(map[string]*domain.Route),
	}
}

func (r *MemoryRepo) Save(ctx context.Context, route *domain.Route) (*domain.Route, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[route.ID] = route

	return route, nil
}

func (r *MemoryRepo) FindByID(ctx context.Context, routeID string) (*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	route, exists := r.routes[routeID]
	if !exists {
		err := domain.NewError(domain.ErrorRouteNotFound, "route not found", nil).
			WithMetadata(domain.RouteID, route)

		return nil, err
	}

	return route, nil
}
