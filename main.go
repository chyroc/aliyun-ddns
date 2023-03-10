package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/chyroc/aliyun-ddns/internal/command"
)

// https://help.aliyun.com/document_detail/124923.html

func main() {
	app := &cli.App{
		Name: "aliyun-ddns",
		Commands: []*cli.Command{
			command.Set(),
			command.Get(),
			command.UpdateSet(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
