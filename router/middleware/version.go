package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/target/jenkins-docker-api/version"
)

// Version is a middleware function that injects the Jenkins API version
// information into the middleware chain so it will be logged. This is
// intended for debugging and troubleshooting.
func Version(c *gin.Context) {
	apiVer := version.Version
	if gin.Mode() == "debug" {
		c.Request.Header.Set("J-Api-Version", apiVer.String())
	} else { // in prod we don't want the build number metadata
		apiVer.Metadata = ""
		c.Request.Header.Set("J-Api-Version", apiVer.String())
	}
}
