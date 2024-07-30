package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prajjwal-w/cetec_golang_practical/routes"
)

func main() {
	r := gin.Default()

	r.Use(gin.Logger())
	routes.Routes(r)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while loading env: %v", err)
	}

	port := os.Getenv("PORT")
	err = r.Run(":" + port)
	if err != nil {
		log.Fatal("error while starting server: ", err)
	}

}
