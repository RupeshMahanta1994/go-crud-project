package main

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func main() {
	app := fiber.New()

	//app.Use(recover.New())

	//Check panic mode at / endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		panic("Panic at disco")
		return c.SendString("Testing endpoint")
	})
	//Implementation of middleware
	app.Use("/orders/code/:orderCode", func(c *fiber.Ctx) error {
		var correlationId = c.Get("x-correlationid")
		if correlationId == "" {
			return c.Status(http.StatusBadRequest).JSON("x-correlationid must mandatory")
		}
		_, err := uuid.Parse(correlationId)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON("Correlationn id must be guid")
		}

		c.Locals("correaltionId", correlationId)
		return c.Next()

	})

	app.Get("/orders/code/:orderCode", func(c *fiber.Ctx) error {
		fmt.Println("Your correlatinId is", c.Locals("correaltionId"))
		return c.SendString("Order code that you sent is very good: " + c.Params("orderCode"))
	})

	app.Get("/users/name/:name", func(c *fiber.Ctx) error {
		if c.Params("name") != "" {
			return c.SendString("Hello " + c.Params("name"))
		}
		return c.SendString("Where is mfk`")
	})

	//Define struct for POST request
	type CreateOrderRequest struct {
		ShipmentNumber string `json:"shipmentNumber"`
	}
	app.Post("/orders", func(c *fiber.Ctx) error {
		var request CreateOrderRequest

		err := c.BodyParser(&request)
		if err != nil {
			return err
		}
		return c.Status(http.StatusCreated).JSON(request)

	})
	//defining struct for the user
	type User struct {
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		Phone    float32 `json:"phone"`
		Password string  `json:"password"`
	}
	app.Post("/register", func(c *fiber.Ctx) error {
		var request User
		c.BodyParser(&request)

		return c.Status(http.StatusAccepted).JSON(request)
	})

	app.Listen(":5000")

}
