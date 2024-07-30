package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prajjwal-w/cetec_golang_practical/routes"
)

func main() {

	//creating a router
	r := gin.New()

	//using the gin's inbuild logger
	r.Use(gin.Logger())
	routes.Routes(r)

	//loading the env file            (commiting the env file for the reference)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while loading env: %v", err)
	}

	port := os.Getenv("PORT")

	//starting the server
	err = r.Run(":" + port)
	if err != nil {
		log.Fatal("error while starting server: ", err)
	}

}
