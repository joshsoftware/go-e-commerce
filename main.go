package main

// @APITitle Main
// @APIDescription Main API for Microservices in Go!

import (
	"fmt"
	"joshsoftware/go-e-commerce/config"
	"joshsoftware/go-e-commerce/db"
	"joshsoftware/go-e-commerce/service"
	"os"
	"strconv"

	"github.com/rs/cors"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	config.Load()

	cliApp := cli.NewApp()
	cliApp.Name = config.AppName()
	cliApp.Version = "1.0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start server",
			Action: func(c *cli.Context) error {
				return startApp()
			},
		},
		{
			Name:  "create_migration",
			Usage: "create migration file",
			Action: func(c *cli.Context) error {
				return db.CreateMigrationFile(c.Args().Get(0))
			},
		},
		{
			Name:  "migrate",
			Usage: "run db migrations",
			Action: func(c *cli.Context) error {
				return db.RunMigrations()
			},
		},
		{
			Name:  "rollback",
			Usage: "rollback migrations",
			Action: func(c *cli.Context) error {
				return db.RollbackMigrations(c.Args().Get(0))
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func startApp() (err error) {
	store, err := db.Init()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Database init failed")
		return
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		Debug:            true,
	})

	deps := service.Dependencies{
		Store: store,
	}

	// mux router
	router := service.InitRouter(deps)

	// init web server
	server := negroni.Classic()
	server.Use(c)
	server.UseHandler(router)

	port := config.AppPort() // This can be changed to the service port number via environment variable.
	addr := fmt.Sprintf(":%s", strconv.Itoa(port))

	server.Run(addr)
	return
}
