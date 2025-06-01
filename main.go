package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default values.")
	}
	
	app := fiber.New()

	// ดึงค่าจาก ENV หรือ fallback เป็นค่า default
	ok := os.Getenv("OK")
	if ok == "" {
		ok = "000"
	}

	ko := os.Getenv("KO")
	if ko == "" {
		ko = "111"
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(ok + " " + ko)
	})

	app.Listen(":" + "80")
}