package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/url"
)

var scope *SQLTransactionScope

func main() {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "root", "127.0.0.1", "3307", "transaction_table")

	val := url.Values{}
	val.Add("parseTime", "true")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatalf("cannot open mysql connection")
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("error connect to mysql")
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatalf("error closing mysql")
		}
	}()

	scope = &SQLTransactionScope{db: dbConn}

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	e.GET("/", handleGET)
	e.POST("/", handlePOST)

	err = e.Start(fmt.Sprintf(":%v", 8080))
	if err != nil {
		log.Fatal(err.Error())
	}
}
