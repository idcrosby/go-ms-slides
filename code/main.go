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

// START OMIT
func main() {
	errc := make(chan error)
	ctx := context.Background()

	db, err := sqlx.Open("mysql", *dsn)
	if err != nil {
		logger.Fatal("Database error: ", err)
	}
	defer db.Close()

	service := catalogue.NewCatalogueService(db)
	endpoints := catalogue.MakeEndpoints(service)

	go func() {
		handler := catalogue.MakeHTTPHandler(ctx, endpoints, *images)
		errc <- http.ListenAndServe(":"+*port, handler)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
}
// END OMIT