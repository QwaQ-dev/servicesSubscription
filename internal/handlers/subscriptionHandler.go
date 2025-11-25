package handlers

import (
	"log/slog"
	"strconv"

	"github.com/QwaQ-dev/servicesSubscription/internal/services"
	"github.com/QwaQ-dev/servicesSubscription/internal/structures"
	"github.com/QwaQ-dev/servicesSubscription/pkg/sl"
	"github.com/gofiber/fiber/v2"
)

type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
	log                 *slog.Logger
}

func NewSubsriptionHandler(
	subscriptionService *services.SubscriptionService,
	log *slog.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		log:                 log,
	}
}

// CreateSubscription godoc
// @Summary Create Subscription
// @Description Creating new subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body structures.Subscription true "Subscription data"
// @Success 200 {object} map[string]interface{} "message + id"
// @Failure 400 {object} structures.ErrorResponse "Invalid subscription format"
// @Failure 500 {object} structures.ErrorResponse "Error"
// @Router /subscription/ [post]
func (h *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	const op = "handlers.subscriptionHandler.CreateSubscription"
	log := h.log.With("op", op)

	subscription := new(structures.Subscription)

	if err := c.BodyParser(subscription); err != nil {
		log.Error("Invalid subscription fromat", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid subscription format",
		})
	}

	id, err := h.subscriptionService.CreateSub(subscription)
	if err != nil {
		log.Error("Failed to create subsciption", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Subscription created successfully",
		"id":      id,
	})
}

// GetAllSubscriptions godoc
// @Summary Get All subscriptions
// @Description List of all subscriptions
// @Tags Subscriptions
// @Produce json
// @Success 200 {object} map[string][]structures.Subscription
// @Failure 500 {object} structures.ErrorResponse
// @Router /subscription/ [get]
func (h *SubscriptionHandler) GetAllSubscriptions(c *fiber.Ctx) error {
	const op = "handlers.subscriptionHandler.GetAllSubscriptions"
	log := h.log.With("op", op)

	subscriptions, err := h.subscriptionService.GetAllSubs()
	if err != nil {
		log.Error("Failed to get all subscriptions", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"subscriptions": subscriptions,
	})
}

// GetOneSubscription godoc
// @Summary Get one subscription by ID
// @Description returns subscription by ID
// @Tags Subscriptions
// @Produce json
// @Param id path int true "subscription ID"
// @Success 200 {object} structures.Subscription
// @Failure 400 {object} structures.ErrorResponse
// @Failure 404 {object} structures.ErrorResponse
// @Router /subscription/{id} [get]
func (h *SubscriptionHandler) GetOneSubscription(c *fiber.Ctx) error {
	const op = "handlers.subscriptionHandler.GetOneSubscription"
	log := h.log.With("op", op)

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid article id", slog.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	subscription, err := h.subscriptionService.GetSubById(id)
	if err != nil {
		log.Error("Subscription not found", slog.Int("id", id), slog.Any("err", err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subscription not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"Subscription": subscription,
	})
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "subscription ID"
// @Param subscription body structures.Subscription true "subscription data"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} structures.ErrorResponse
// @Failure 500 {object} structures.ErrorResponse
// @Router /subscription/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *fiber.Ctx) error {
	const op = "handlers.subscriptionHandler.UpdateSubscription"
	log := h.log.With("op", op)

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid article id", slog.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var subscription structures.Subscription
	if err := c.BodyParser(&subscription); err != nil {
		log.Error("Failed to parse subscription body", slog.Any("err", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscription format",
		})
	}

	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required for update",
		})
	}

	if err := h.subscriptionService.UpdateSub(&subscription, id); err != nil {
		log.Error("Failed to update subscription", slog.Any("err", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update subscription",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Subscription has been updated",
	})
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags Subscriptions
// @Produce json
// @Param id path int true "subscription ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} structures.ErrorResponse
// @Failure 500 {object} structures.ErrorResponse
// @Router /subscription/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *fiber.Ctx) error {
	const op = "handlers.subscriptionHandler.DeleteSubscription"
	log := h.log.With("op", op)

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("Invalid article id", slog.String("id", idStr))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.subscriptionService.DeleteSub(id); err != nil {
		log.Error("Failed to delete article", slog.Int("id", id), slog.Any("err", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete subscription",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Subscription has been deleted",
	})
}

// GetSumm godoc
// @Summary Get summ of subscriptions prices
// @Description Returns the sum of all subscriptions filtered by date, user and service
// @Tags Sum
// @Accept json
// @Produce json
// @Param counting body structures.Counting true "Filters"
// @Success 200 {object} map[string]float64 "total"
// @Failure 400 {object} structures.ErrorResponse
// @Failure 500 {object} structures.ErrorResponse
// @Router /summ/ [get]
func (h *SubscriptionHandler) GetSumm(c *fiber.Ctx) error {
	const op = "handlers.subscriptionHandler.GetSumm"
	log := h.log.With("op", op)

	var data structures.Counting
	if err := c.BodyParser(&data); err != nil {
		log.Error("Failed to parse counting body", sl.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	total, err := h.subscriptionService.Counting(&data)
	if err != nil {
		log.Error("Failed to get sum", sl.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get sum",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"total": total,
	})
}
