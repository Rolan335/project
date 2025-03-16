package middleware

import (
	"strconv"

	"github.com/Rolan335/project/internal/metric"
	"github.com/gofiber/fiber/v2"
)

func Metric(c *fiber.Ctx) error {
	//Копируем метод для правильного отображения в метрике
	method := string([]byte(c.Method()))

	err := c.Next()

	metric.RequestsCounter.WithLabelValues(method, strconv.Itoa(c.Response().StatusCode())).Inc()

	return err
}
