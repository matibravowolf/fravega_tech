package app

import (
	"github.com/uMakeMeCrazy/fravega_tech/internal/core/domain"

	"github.com/uMakeMeCrazy/fravega_tech/pkg/middleware"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Start() {
	// Init logger
	baseLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// Init dependencies
	dep := initDependencies()

	// Configure router
	r := gin.New()
	r.Use(middleware.WithLogger(baseLogger))
	r.Use(domain.ErrorHandler())

	addRoutes(r, dep)

	// Run
	_ = r.Run(":8080")
}
