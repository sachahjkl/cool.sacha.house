package common

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func RedirectHTMX(c *fiber.Ctx, routeName string, params fiber.Map) error {

	route, err := c.GetRouteURL(routeName, params)
	if err != nil {
		return err
	}

	c.Append("HX-Redirect", route)
	return c.SendStatus(fiber.StatusNoContent)
}

func ClearUser(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name: USER_COOKIE_KEY,
		// expires in the past
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})
}
