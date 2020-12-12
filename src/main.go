package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"context"
	"strings"
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

	app.Static("/", "/public")

	app.Get("*", func(c *fiber.Ctx) error {
		c.SendFile("/public/index.html")
		return nil
	})

	app.Listen(":80")
}

func api(app *fiber.App) fiber.Router {
	api := app.Group("/api", auth)

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

func auth(c *fiber.Ctx) error {
	// Parse creadentials from the request header
	// Decode credentials, convert to `email` and `password`
	// Find user in database with specified credentials
	// Attach user info (`id`, `name`, `login`) to request and pass to Next
	// If any of the steps failed return Unauthrized error

	h := fmt.Sprintf("%s", c.Request().Header.Peek("Authorization"))	

	switch auth := strings.Fields(h); strings.ToLower(auth[0]) {
	case "basic":
		c.Next()
		return nil
	}

	return fiber.ErrUnauthorized
}