package gin

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes is a helper function to demonstrate route setup
// You can customize this based on your needs
func (e *Engine) SetupRoutes(routes func(router *gin.Engine)) {
	routes(e.router)
}
