package routeshdl

type createRouteRequest struct {
	Vehicle string `json:"vehicle" binding:"required"`
	Driver  string `json:"driver" binding:"required"`
}
