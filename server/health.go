package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /health health getHealth
//
// Check if the Jenkins API is available
//
// ---
// x-success_http_code: '200'
// x-response_time_sla: 100ms
// x-expected_tps: '1'
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful 'ping' of Jenkins API
//     schema:
//       type: string
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
