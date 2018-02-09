package router

import (
	"github.com/gin-gonic/gin"
	"github.com/target/jenkins-docker-api/router/middleware/header"
	"github.com/target/jenkins-docker-api/server"
)

// Load is a server function that returns the engine for processing web requests
// on the host it's running on
func Load(options ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(header.Version)
	r.Use(header.NoCache)
	r.Use(header.Options)
	r.Use(header.Secure)

	r.Use(options...)

	base := r.Group("/api/v1")
	{
		base.GET("/health", server.Health)
		base.GET("/jenkins", server.GetAll)
		base.PUT("/jenkins/restart/:master", server.Restart)
		base.PUT("/jenkins/update/:master", server.Update)
		base.PUT("/jenkins/shutdown/:master", server.Shutdown)
		base.PUT("/jenkins/admin/update_all", server.UpdateAll)
	}

	return r
}
