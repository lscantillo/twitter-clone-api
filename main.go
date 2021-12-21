package main

import (
	"log"

	"github.com/lscantillo/twitter-clone-api/config"
	"github.com/lscantillo/twitter-clone-api/db"
	"github.com/lscantillo/twitter-clone-api/routers"
)

func main() {
	if db.CheckDB() == 0 {
		log.Fatal("Error while connecting to the database")
		return
	}
	a := routers.App{}
	a.Initialize()
	a.Run(config.GetVariables("PORT"))
}
