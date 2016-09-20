package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/microservices-demo/catalogue"
	"golang.org/x/net/context"
	"path/filepath"
)

func main() {

	errc := make(chan error)
	ctx := context.Background()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	db, err := sqlx.Open("mysql", *dsn)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	defer db.Close()

	var service catalogue.Service
	{
		service = catalogue.NewCatalogueService(db, logger)
		service = catalogue.LoggingMiddleware(logger)(service)
	}

	endpoints := catalogue.MakeEndpoints(service)

	go func() {
		logger.Log("transport", "HTTP", "port", *port)
		handler := catalogue.MakeHTTPHandler(ctx, endpoints, *images, logger)
		errc <- http.ListenAndServe(":"+*port, handler)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
