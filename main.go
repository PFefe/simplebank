package main

import (
	"database/sql"
	"github.com/PFefe/simplebank/api"
	db "github.com/PFefe/simplebank/db/sqlc"
	"github.com/PFefe/simplebank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(
			"cannot load config:",
			err,
		)
	}
	conn, err := sql.Open(
		config.DBDriver,
		config.DBSource,
	)
	if err != nil {
		log.Fatal(
			"cannot connect to db:",
			err,
		)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal(
			"cannot start server:",
			err,
		)
	}

}
