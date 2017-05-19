package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"net"
	"bufio"
	"strings"
	"flag"
)

var datastoreAddr string

func getUser(c echo.Context) error {
	conn, err := net.Dial("tcp", datastoreAddr)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error connecting to the datastore")
	}
	defer conn.Close()

	name := c.Param("name")
	fmt.Fprintf(conn, "%s\n", name)

	greeting, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading from the datastore")
	}
	return c.String(http.StatusOK, fmt.Sprintf("Hi, %s!", strings.TrimSpace(greeting)))
}

const (
	DATASTORE_HOST = "localhost"
	DATASTORE_PORT = "3333"
)

func main() {
	datastoreHost := flag.String("datastore-host", DATASTORE_HOST, "IP address of the datastore")
	datastorePort := flag.String("datastore-port", DATASTORE_PORT, "port of the datastore")

	flag.Parse()

	datastoreAddr = *datastoreHost+":"+*datastorePort

	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users/:name", getUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
