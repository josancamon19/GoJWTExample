package main

import (
	"ConchaAPI/auth"
	"ConchaAPI/database"
	"ConchaAPI/models"
	"ConchaAPI/views"
	"fmt"
	"github.com/gofiber/fiber"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func initDatabase() {
	var err error
	dsn := "postgres://ltfukhqavcelvb:4dcb1966901321044ea8916ebb86520064a9c5426b54687511dec5eb39cb765d@ec2-3-95-87-221.compute-1.amazonaws.com:5432/dakstsaagpqe2g"
	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = database.DBConn.AutoMigrate(&models.User{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Database migrated")
}

func main() {
	initDatabase()
	authHandler := &auth.Handler{}

	app := fiber.New(fiber.Config{})
	api := app.Group("/api")
	apiV1 := api.Group("/v1")
	apiV1.Post("/auth", authHandler.Auth)
	apiV1.Post("/refresh", authHandler.RefreshToken)
	apiV1.Post("/users", views.CreateUserEndpoint)
	apiV1.Get("/users", views.GetUsersEndpoint)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
