package main

import (
	"github.com/etg-dev/restApi/configs"
	"github.com/etg-dev/restApi/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//! seed program
	//err := seeders.DrainDB()
	// err := seeders.InjectDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//!

	configs.ConnectDB()

	routes.UserRoute(router)
	routes.PostRoute(router)

	router.Run("localhost:6000")
}
