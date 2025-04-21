package routeshdl

import (
	"context"
	"net/http"

	"github.com/uMakeMeCrazy/fravega_tech/internal/core/domain"

	"github.com/uMakeMeCrazy/fravega_tech/pkg/logger"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type RoutesService interface {
	CreateRoutes(ctx context.Context, driver string, vehicle string) (*domain.Route, error)
	AddPurchaseToRoute(ctx context.Context, routeID string, purchaseID string) error
	GetRoute(ctx context.Context, routeID string) (*domain.Route, error)
	PurchaseDeliveredNotification(ctx context.Context, routeID string, purchaseID string) error
}

type RoutesHdl struct {
	routesSrv RoutesService
}

func NewRoutesHdl(routesSrv RoutesService) *RoutesHdl {
	return &RoutesHdl{
		routesSrv: routesSrv,
	}
}

func (r *RoutesHdl) CreateRoute(c *gin.Context) {
	ctx := c.Request.Context()

	var req createRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		customErr := domain.NewError(domain.ErrorBadRequest, "invalid request", err)
		_ = c.Error(customErr)
		return
	}

	logger.Info(ctx, "starting of CreateRoutes()", zap.Any("request", req))
	res, err := r.routesSrv.CreateRoutes(ctx, req.Driver, req.Vehicle)
	if err != nil {
		logger.Error(ctx, "failed to CreateRoutes()", zap.Error(err))
		_ = c.Error(err)
		return
	}
	logger.Info(ctx, "ending of CreateRoutes()", zap.String(domain.RouteID, res.ID))

	c.JSON(http.StatusCreated, res)
}

func (r *RoutesHdl) AddPurchaseToRoute(c *gin.Context) {
	ctx := c.Request.Context()

	routeID := c.Param(domain.RouteID)
	purchaseID := c.Param(domain.PurchaseID)

	logger.Info(ctx, "starting of AddPurchaseToRoute()", zap.String(domain.RouteID, routeID), zap.String(domain.PurchaseID, purchaseID))
	err := r.routesSrv.AddPurchaseToRoute(ctx, routeID, purchaseID)
	if err != nil {
		logger.Error(ctx, "failed of AddPurchaseToRoute()", zap.String(domain.RouteID, routeID), zap.String(domain.PurchaseID, purchaseID), zap.Error(err))
		_ = c.Error(err)
		return
	}
	logger.Info(ctx, "ending of AddPurchaseToRoute()", zap.String(domain.RouteID, routeID), zap.String(domain.PurchaseID, purchaseID))

	c.Status(http.StatusOK)
}

func (r *RoutesHdl) GetRoute(c *gin.Context) {
	ctx := c.Request.Context()
	routeID := c.Param(domain.RouteID)

	logger.Info(ctx, "starting of GetRoute()", zap.String(domain.RouteID, routeID))
	route, err := r.routesSrv.GetRoute(ctx, routeID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Info(ctx, "ending of GetRoute()", zap.String(domain.RouteID, route.ID))

	c.JSON(http.StatusOK, route)
}

func (r *RoutesHdl) PurchaseDeliveredNotification(c *gin.Context) {
	ctx := c.Request.Context()

	routeID := c.Param(domain.RouteID)
	purchaseID := c.Param(domain.PurchaseID)

	logger.Info(ctx, "starting of PurchaseDeliveredNotification()", zap.String(domain.RouteID, routeID), zap.String(domain.PurchaseID, purchaseID))
	err := r.routesSrv.PurchaseDeliveredNotification(ctx, routeID, purchaseID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Info(ctx, "ending of PurchaseDeliveredNotification()", zap.String(domain.RouteID, routeID), zap.String(domain.PurchaseID, purchaseID))

	c.JSON(http.StatusOK, nil)
}
