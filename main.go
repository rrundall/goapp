package main

import (
	"flag"
	"goapp/pkg/api"
	"goapp/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//	@title			Book Library API
//	@version		1.0
//	@description	This is a simple book library server.

//	@host		localhost:8080
//	@BasePath	/v1

// Main function to start up the server and all services.
func main() {
	// Enter --addr/--debug flag options if you don't want to use default setup
	// By default the server will start on localhost:8080 and debug level is set to false
	address := flag.String("addr", ":8080", "host address")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()
	d := db.OpenSqliteStorage()
	defer d.CloseDB()
	router := gin.New()
	router.Use(gin.Recovery())
	server := api.GetServer(*address, router, d)
	server.StartLogging(debug)
	log.Info().Msgf("Server is running port -> %s", *address)
	log.Fatal().Err(server.StartServer()).Msg("fail to start server")
}
