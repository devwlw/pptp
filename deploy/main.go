package main

import (
	"log"
	"mail/deploy/autojob"
	"mail/deploy/context"
	"mail/deploy/router"
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
			Usage:  "start a api server",
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
	log.Println("server start at port 9100")
	go func() {
		log.Println("auto job is running")
		autojob.AutoJob{}.Do()
	}()
	err := http.ListenAndServe(":9100", handlers.LoggingHandler(os.Stdout, routers))
	if err != nil {
		panic(err)
	}

}
