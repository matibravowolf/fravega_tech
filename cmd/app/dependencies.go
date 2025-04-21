package app

import (
	"github.com/uMakeMeCrazy/fravega_tech/cmd/handlers/pinghdl"
	"github.com/uMakeMeCrazy/fravega_tech/cmd/handlers/routeshdl"
	"github.com/uMakeMeCrazy/fravega_tech/internal/core/services"
	"github.com/uMakeMeCrazy/fravega_tech/internal/repositories"
)

type Dependencies struct {
	// Handlers
	PingHdl   *pinghdl.PingHdl
	RoutesHdl *routeshdl.RoutesHdl

	// Services
	RoutesSrv *services.RoutesSrv

	// Repositories
	routesRepo    *repositories.MemoryRepo
	purchasesRepo *repositories.PurchasesRepo
}

func initDependencies() Dependencies {
	// Repositories
	routesRepo := repositories.NewMemoryRepository()
	purchasesRepo := repositories.NewPurchasesRepo()

	// Services
	routesSrv := services.NewRoutesSrv(routesRepo, purchasesRepo)

	// Handlers
	pingHdl := pinghdl.NewPingHdl()
	routesHdl := routeshdl.NewRoutesHdl(routesSrv)

	return Dependencies{
		PingHdl:   pingHdl,
		RoutesHdl: routesHdl,
	}
}
