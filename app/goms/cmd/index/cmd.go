package index

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/dstarapp/gomorainsc/http"
	"github.com/dstarapp/gomorainsc/indexer"
	"github.com/dstarapp/gomorainsc/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli"
)

var Cmd = cli.Command{
	Name:        "index",
	Usage:       "index -c conffile -b blackfile",
	Description: "mora indexer",
	Action:      runCmd,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Value: "./conf/app.conf",
			Usage: "config file path",
		},
		cli.StringFlag{
			Name:  "b",
			Value: "./conf/black.txt",
			Usage: "black file path",
		},
		cli.IntFlag{
			Name:  "p",
			Value: 8301,
			Usage: "server listen port",
		},
	},
}

func runCmd(ctx *cli.Context) {
	cfg, err := get_local_config(ctx.String("c"), ctx.String("b"))
	if err != nil {
		logs.Error(err)
		return
	}

	done := make(chan bool)
	if err := http.StartServer(cfg, ctx.Int("p"), done); err != nil {
		logs.Error(err)
	}
	<-done
}

func get_local_config(conffile, blackfile string) (*indexer.Config, error) {
	data, err := os.ReadFile(conffile)
	if err != nil {
		return nil, err
	}
	cfg := indexer.Config{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	blackids, err := utils.ReadFileByLine(blackfile)
	if err != nil {
		return &cfg, nil
	}

	cfg.BlackIDs = blackids

	return &cfg, nil
}
