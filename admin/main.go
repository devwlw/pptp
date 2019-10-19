package main

import (
	"log"
	"mail/admin/context"
	"mail/admin/router"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/urfave/cli"
)

func main() {
	NewCliApp().Run(os.Args)
}

func NewCliApp() *cli.App {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "start a api sersver",
			Action: serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Usage: "server port",
				},
				cli.StringFlag{
					Name:  "config,c",
					Value: "./config.json",
					Usage: "the config file path",
				},
			},
		},
	}
	return app
}

func serve(c *cli.Context) {
	confPath := c.String("config")
	context.InitContext(confPath)
	routers := router.NewRouter()
	log.Println("admin server start at port 80")
	err := http.ListenAndServe(":80", handlers.LoggingHandler(os.Stdout, routers))
	if err != nil {
		panic(err)
	}
	/*	confPath := c.String("config")
		config, err := conf.NewConfig(confPath)
		if err != nil {
			panic(err)
		}
		util.InitRegistry(config)
		conf.ConfigInstance = config
		//insertFake()
		//insertFakeCashout()

		flagPort := config.Port
		routers := router.NewRouter()
		log.Printf("server start at: %s", flagPort)
		http.Handle("/vendors/", http.StripPrefix("/vendors/", http.FileServer(http.Dir("templates/vendors"))))
		err = http.ListenAndServe(":"+flagPort, handlers.LoggingHandler(os.Stdout, routers))
		if err != nil {
			log.Fatal(err)
		}*/
}
