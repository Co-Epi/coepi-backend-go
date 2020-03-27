package main

import (
	"flag"
	"fmt"

	"github.com/Co-Epi/coepi-backend-go/backend"
	"github.com/Co-Epi/coepi-backend-go/server"
)

const (
	version = "0.2"
)

func main() {
	port := flag.Uint("port", uint(server.DefaultPort), "port cen is listening on")
	mysqlConnectionString := flag.String("conn", backend.DefaultConnString, "MySQL Connection String")

	_, err := server.NewServer(uint16(*port), *mysqlConnectionString)
	if err != nil {
		panic(err)
	}
	fmt.Printf("CEN Server v%s - Listening on port %d...\n", version, *port)
	for {
	}
}
