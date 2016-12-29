package main

import (
	"context"
	"net/http"

	"github.com/FarmRadioHangar/devices/db"
	"github.com/FarmRadioHangar/devices/events"
	"github.com/FarmRadioHangar/devices/udev"
	"github.com/FarmRadioHangar/devices/web"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Usage = "Sreams realtime events about devices (Dongles)"
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Starts a server that listens to udev vents",
			Action:  Server,
		},
	}
}

func Server(cxt *cli.Context) error {
	s := events.NewStream(1000)
	ql, err := db.DB()
	if err != nil {
		return err
	}
	m := udev.New(s)
	ctx, cancel := context.WithCancel(context.Background())
	go m.Run(ctx)

	defer cancel()
	w := web.New(ql, s)
	return http.ListenAndServe(":1000", w)
}
