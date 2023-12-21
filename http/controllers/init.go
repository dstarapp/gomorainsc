package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func InitRouter(app *iris.Application) {
	mvc.Configure(app.Party("/v1"),
		myMVC,
	)
}

func myMVC(app *mvc.Application) {
	api := app.Party("/api")

	api.Party("/ft").Handle(new(FtController))
	api.Party("/ftitem").Handle(new(FtItemController))
}
