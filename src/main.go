package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"context"
	"fmt"
	"os"
)

var conn *pgx.Conn

func main() {
	// Init DB connection
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Start web server
	app := fiber.New()

	api(app)

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	app.Listen(":80")
}

func api(app *fiber.App) fiber.Router {
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		var id int64
		var name, email string

		err := conn.QueryRow(context.Background(), "SELECT id, email, name FROM users").Scan(&id, &email, &name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			return err
		}

		c.JSON(fiber.Map {
			"id": id,
			"name": name,
			"email": email,
		})

		return nil
	})

	return api
}