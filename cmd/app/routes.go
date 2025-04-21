package app

import (
	"github.com/gin-gonic/gin"
)

func addRoutes(r *gin.Engine, dep Dependencies) {
	// Health
	r.GET("/ping", dep.PingHdl.Ping)

	// Routes v1
	v1 := r.Group("/api/v1")
	{
		routes := v1.Group("/routes")
		{
			routes.POST("/", dep.RoutesHdl.CreateRoute)
			routes.GET("/:route_id", dep.RoutesHdl.GetRoute)

			// Purchases
			routes.PUT("/:route_id/purchases/:purchase_id", dep.RoutesHdl.AddPurchaseToRoute)

			// Email notification
			routes.POST("/:route_id/purchases/:purchase_id/notification", dep.RoutesHdl.PurchaseDeliveredNotification)
		}
	}
}
