package pinghdl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHdl struct{}

func NewPingHdl() *PingHdl {
	return &PingHdl{}
}

func (p *PingHdl) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
