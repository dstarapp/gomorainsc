package http

import (
	"context"
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/dstarapp/gomorainsc/http/controllers"
	"github.com/dstarapp/gomorainsc/http/services"
	"github.com/dstarapp/gomorainsc/indexer"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	requestLogger "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func stop_server() {
	services.StopIndexer()
	logs.Info("server stop success")

}

func StartServer(cfg *indexer.Config, port int, exit chan bool) error {
	if err := services.StartIndexer(cfg); err != nil {
		return nil
	}
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	InitServer(addr, func() {
		stop_server()
		exit <- true
	})
	return nil
}

func InitServer(addr string, shutdown func()) error {
	app := iris.New()
	cfg := requestLogger.DefaultConfig()
	cfg.Query = true

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	// app.Use(crs.Handler())
	app.UseRouter(crs)
	app.Use(recover.New())
	app.Use(requestLogger.New(cfg))
	app.Use(iris.Compression)
	// //app.StaticWeb("/", "statistic")
	controllers.InitRouter(app)

	if shutdown != nil {
		iris.RegisterOnInterrupt(func() {
			app.Shutdown(context.TODO())
			shutdown()
		})
	}

	err := app.Run(iris.Addr(addr),
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
	)
	return err
}
