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
	store, err := db.Init() //pg.go - datatype will be &pgStore - *sql.DB
	if err != nil {
		logger.WithField("err", err.Error()).Error("Database init failed")
		return
	}

	deps := service.Dependencies{
		Store: store,
		// here value assigned to Store which is of type interface Storer ? and value of type &pgStore - *sql.DB
		// 1st how are we able to create object of Store which is a interface
		// 2nd how are we assigning it to variable of another type
	}

	// mux router
	router := service.InitRouter(deps) // init router - return mux.NewRouter with all hanldefuncs
	//in router.go why have in some places we have used hanldefunc and some handle

	// init web server
	server := negroni.Classic()
	server.UseHandler(router)

	port := config.AppPort() // This can be changed to the service port number via environment variable.
	addr := fmt.Sprintf(":%s", strconv.Itoa(port))

	server.Run(addr)
	return
}
