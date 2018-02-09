// Package classification Jenkins API
//
// Self Service API built by Target to manage Jenkins master
//
// Terms Of Service:
//
// http://swagger.io/terms/
//
//     Version: 1.0.0
//     Host: yourhost.com
//     Schemes: https
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"net/http"
	"time"

	tomb "gopkg.in/tomb.v2"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/target/jenkins-docker-api/router/middleware/logger"

	"github.com/target/jenkins-docker-api/router"
	"github.com/target/jenkins-docker-api/router/middleware"

	log "github.com/Sirupsen/logrus"
)

func server(c *cli.Context) error {
	// debug level if requested by user
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.Debug("Starting server")

	log.SetFormatter(&log.JSONFormatter{})
	router := router.Load(
		middleware.Version,
		logger.LogHandler(log.StandardLogger(), time.RFC3339, true),
	)

	log.Debug("Server started")

	var tomb tomb.Tomb

	// start http server
	tomb.Go(func() error {
		srv := &http.Server{Addr: c.String("server-port"), Handler: router}

		go func() {
			err := srv.ListenAndServe()
			if err != nil {
				tomb.Kill(err)
			}
		}()

		log.Debug("Started HTTP")

		for {
			select {
			case <-tomb.Dying():
				log.Debug("Stopping HTTP")
				return srv.Shutdown(nil)
			}
		}
	})

	// Wait for stuff and watch for errors
	tomb.Wait()
	return tomb.Err()
}
