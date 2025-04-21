package services

import (
	"context"

	"github.com/uMakeMeCrazy/fravega_tech/pkg/logger"
	"go.uber.org/zap"

	"github.com/uMakeMeCrazy/fravega_tech/internal/core/domain"
)

type RoutesRepo interface {
	Save(ctx context.Context, route *domain.Route) (*domain.Route, error)
	FindByID(ctx context.Context, routeID string) (*domain.Route, error)
}

type PurchasesRepo interface {
	FindByID(ctx context.Context, purchaseID string) (*domain.Purchase, error)
	SendEmailNotification(ctx context.Context, purchaseID string) error
}

type RoutesSrv struct {
	routesRepo    RoutesRepo
	purchasesRepo PurchasesRepo
}

func NewRoutesSrv(routesRepo RoutesRepo, purchasesRepo PurchasesRepo) *RoutesSrv {
	return &RoutesSrv{
		routesRepo:    routesRepo,
		purchasesRepo: purchasesRepo,
	}
}

func (r *RoutesSrv) CreateRoutes(ctx context.Context, driver string, vehicle string) (*domain.Route, error) {
	route := domain.NewRoute(vehicle, driver)
	savedRoute, err := r.routesRepo.Save(ctx, route)

	return savedRoute, err
}

func (r *RoutesSrv) AddPurchaseToRoute(ctx context.Context, routeID string, purchaseID string) error {
	route, err := r.GetRoute(ctx, routeID)
	if err != nil {
		return err
	}

	if isPending := route.IsPendingRoute(); !isPending {
		return domain.NewError(domain.ErrorPurchaseNotInPendingStatus, "can't add purchase to route, is not in pending status", nil)
	}

	purchase, err := r.purchasesRepo.FindByID(ctx, purchaseID)
	if err != nil {
		return err
	}

	err = route.AddPurchase(purchase)
	if err != nil {
		return err
	}

	_, err = r.routesRepo.Save(ctx, route)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoutesSrv) GetRoute(ctx context.Context, routeID string) (*domain.Route, error) {
	route, err := r.routesRepo.FindByID(ctx, routeID)
	logger.Info(ctx, "success call to routesRepo FindByID()", zap.String(domain.RouteID, routeID))

	return route, err
}

func (r *RoutesSrv) PurchaseDeliveredNotification(ctx context.Context, routeID string, purchaseID string) error {
	route, err := r.routesRepo.FindByID(ctx, routeID)
	if err != nil {
		return err
	}

	err = route.MarkPurchaseAsDelivered(purchaseID)
	if err != nil {
		return err
	}

	_, err = r.routesRepo.Save(ctx, route)
	if err != nil {
		return err
	}

	err = r.purchasesRepo.SendEmailNotification(ctx, purchaseID)
	if err != nil {
		return err
	}

	return nil
}
