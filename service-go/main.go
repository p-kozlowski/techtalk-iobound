package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"net"
	"bufio"
	"strings"
)

func getUser(c echo.Context) error {
	conn, err := net.Dial("tcp", "127.0.0.1:3333")
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

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users/:name", getUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
