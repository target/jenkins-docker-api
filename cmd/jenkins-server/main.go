package main

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "jenkins-server"
	app.Action = server

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Debug level logging",
			EnvVar: "JENKINS_DEBUG",
		},
		cli.StringFlag{
			Name:   "gin-mode",
			Usage:  "Gin Mode",
			EnvVar: "GIN_MODE",
		},
		cli.StringFlag{
			Name:   "server-port",
			Usage:  "Jenkins API listen port",
			EnvVar: "JENKINS_API_PORT",
			Value:  ":8080",
		},
		cli.StringFlag{
			Name:   "env",
			Usage:  "Jenkins Environment",
			EnvVar: "JENKINS_ENV",
		},
		cli.StringFlag{
			Name:   "img",
			Usage:  "Jenkins Image",
			EnvVar: "JENKINS_IMAGE",
		},
		cli.StringFlag{
			Name:   "user-config",
			Usage:  "Jenkins User Config Path",
			EnvVar: "JENKINS_USER_CONFIG_PATH",
			Value:  "/jenkins/user-configs/",
		},
		cli.StringFlag{
			Name:   "admin-config",
			Usage:  "Jenkins Admin Config Path",
			EnvVar: "JENKINS_ADMIN_CONFIG_PATH",
			Value:  "/jenkins/secret-configs/",
		},
		cli.StringFlag{
			Name:   "github-api",
			Usage:  "GitHub API",
			EnvVar: "GITHUB_API_URL",
		},
		cli.StringFlag{
			Name:   "admin-org",
			Usage:  "Admin Organization",
			EnvVar: "GITHUB_ADMIN_ORG",
			Value:  "Jenkins",
		},
		cli.StringFlag{
			Name:   "admin-team",
			Usage:  "Admin Team",
			EnvVar: "GITHUB_ADMIN_TEAM",
			Value:  "Admins",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
