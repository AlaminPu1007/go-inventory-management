package main

import (
	"database/sql"
	"log"

	"github.com/alaminpu1007/inventory-system/api"
	db "github.com/alaminpu1007/inventory-system/db/sqlc"
	"github.com/alaminpu1007/inventory-system/utils"

	_ "github.com/lib/pq"
)

func main() {

	// load from app.env
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("DB connection is not possible", err)
	}

	// now create a new store
	store := db.NewStore(conn)

	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	runGinServer(config, store)

}

// This method will run gin server
func runGinServer(config utils.Config, store *db.Store) {
	// using the store, create a new server
	server, err := api.NewServer(config, store)

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
