package main

// import "C"
import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/dstarapp/gomorainsc/app/goms/cmd/index"
	"github.com/urfave/cli"
)

const APP_VER = "1.0"

func main() {

	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	app := cli.NewApp()
	app.Name = "goms"
	app.Usage = "goms help"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		index.Cmd,
	}
	app.Run(os.Args)
}
