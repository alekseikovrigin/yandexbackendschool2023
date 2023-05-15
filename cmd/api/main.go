package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/alekseikovrigin/yandexbackendschool2023/handlers"
	"github.com/alekseikovrigin/yandexbackendschool2023/repository"
	"log"
	"time"
)

func main() {

	// Connect to database
	_, err := repository.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Fiber instance
	app := fiber.New()

	// Use logger
	app.Use(logger.New())

	// Use rate limiter
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
	}))

	// Group Orders related APIs
	orderGroup := app.Group("/orders")

	orderGroup.Get("/", handlers.GetAllOrders)
	orderGroup.Get("/:id", handlers.GetSingleOrder)
	orderGroup.Post("/", handlers.AddNewOrders)

	orderGroup.Post("/complete", handlers.MarkAsComplete) // batch complete
	orderGroup.Post("/assign", handlers.AssignOrders)     // batch complete

	// Group Couriers related APIs
	courierGroup := app.Group("/couriers")

	courierGroup.Get("/", handlers.GetAllCouriers)
	courierGroup.Get("/:id", handlers.GetSingleCourier)
	courierGroup.Post("/", handlers.AddNewCouriers)

	courierGroup.Get("/meta-info/:id", handlers.GetCourierRating) // rating

	errInit := handlers.TypeInsert()
	if errInit != nil {
		return
	}

	err = app.Listen(":3098")
	if err != nil {
		log.Fatal(err)
	}
}
