package middleware

import (
	"log"
	"sachahjkl/htmx_go/pkg/common"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) AuthenticatedOnly(c *fiber.Ctx) error {
	local := c.Locals(common.USER_CONTEXT_KEY)

	if local == nil {
		return c.RedirectToRoute("login", nil)
	}
	userJWT := local.(*jwt.Token)

	user, err := common.UserFromJwt(m.DB, userJWT)

	if err != nil {

		log.Printf("couldn't find user : %v - %v", *userJWT, err)
		// clear cookie, it is invalid since it didn't
		// provide a key to any real user
		c.ClearCookie(common.USER_COOKIE_KEY)
		return c.RedirectToRoute("login", nil)
	}

	// make the user available for next request
	c.Locals(common.USER_KEY, user)
	return c.Next()

}

func (m *Middleware) UnauthenticatedOnly(c *fiber.Ctx) error {
	if c.Locals(common.USER_CONTEXT_KEY) != nil {
		log.Printf("found a user in locals, going back to root")
		return c.RedirectToRoute("root", nil)
	}

	return c.Next()
}

func (m *Middleware) WithUser(c *fiber.Ctx) error {
	local := c.Locals(common.USER_CONTEXT_KEY)

	if local == nil {
		return c.Next()
	}
	userJWT := local.(*jwt.Token)

	user, err := common.UserFromJwt(m.DB, userJWT)

	if err != nil {
		return c.Next()
	}
	c.Locals(common.USER_KEY, user)

	return c.Next()
}
