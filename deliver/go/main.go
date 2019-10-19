package main

import (
	"log"
	"mail/deliver/go/service"
	"os"

	"github.com/urfave/cli"
)

func main() {
	_ = NewCliApp().Run(os.Args)
}

func NewCliApp() *cli.App {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "send",
			Usage:  "",
			Action: action,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type",
					Usage: "163,126",
				},
				cli.StringFlag{
					Name: "receiver",
				},
				cli.StringFlag{
					Name: "title",
				},
				cli.StringFlag{
					Name: "body",
				},
				cli.StringFlag{
					Name: "username",
				},
				cli.StringFlag{
					Name: "password",
				},
				cli.StringFlag{
					Name: "mode",
				},
			},
		},
	}
	return app
}

func action(c *cli.Context) {
	sendType := c.String("type")
	sender, err := service.GetSender(sendType)
	if err != nil {
		sendFailLog(err.Error())
	}
	receiver := c.String("receiver")
	title := c.String("title")
	body := c.String("body")
	username := c.String("username")
	password := c.String("password")
	mode := c.String("mode")
	if receiver == "" {
		sendFailLog("receiver不能为空")
	}
	if title == "" {
		sendFailLog("title不能为空")
	}
	if body == "" {
		sendFailLog("body不能为空")
	}
	if receiver == "" {
		sendFailLog("receiver不能为空")
	}
	if username == "" {
		sendFailLog("username不能为空")
	}
	if password == "" {
		sendFailLog("password不能为空")
	}
	if mode == "" {
		sendFailLog("mode不能为空")
	}
	log.Printf("type:%s,username:%s,password:%s,title:%s,receiver:%s", sendType, username, password, title, receiver)
	err = sender.Send(username, password, receiver, title, body, mode)
	if err != nil {
		sendFailLog(err.Error())
	}
	sendOkLog()
}

func sendOkLog() {
	log.Println("send success")
}

func sendFailLog(msg string) {
	panic("send fail:" + msg)
}
