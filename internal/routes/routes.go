package routes

import (
	"log/slog"

	"github.com/QwaQ-dev/servicesSubscription/internal/handlers"
	swagger "github.com/gofiber/swagger"

	_ "github.com/QwaQ-dev/servicesSubscription/docs"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(
	app *fiber.App,
	log *slog.Logger,
	subscriptionHandler *handlers.SubscriptionHandler,
) {
	v1 := app.Group("/api/v1")

	v1.Get("/swagger/*", swagger.HandlerDefault)

	subscriptionGroup := v1.Group("/subscription")

	subscriptionGroup.Get("/", subscriptionHandler.GetAllSubscriptions)
	subscriptionGroup.Get("/:id", subscriptionHandler.GetOneSubscription)
	subscriptionGroup.Post("/", subscriptionHandler.CreateSubscription)
	subscriptionGroup.Put("/:id", subscriptionHandler.UpdateSubscription)
	subscriptionGroup.Delete("/:id", subscriptionHandler.DeleteSubscription)

	sumGroup := v1.Group("/summ")

	sumGroup.Get("/", subscriptionHandler.GetSumm)
}
