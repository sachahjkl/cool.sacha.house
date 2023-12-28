package handler

import (
	"errors"
	"sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/config"
	"sachahjkl/htmx_go/pkg/middleware"
	"sachahjkl/htmx_go/pkg/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginRequestBody struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	RememberMe bool   `form:"remember-me"`
}

type RegisterRequestBody struct {
	Username        string `form:"username"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"password-confirm"`
}

func RegisterAuthRoutes(app *fiber.App, db *gorm.DB, c *config.Config) {
	h := &handler{
		DB:     db,
		Config: c,
	}

	middleware := middleware.New(db)

	app.Get("/login", middleware.UnauthenticatedOnly, h.LoginPage).Name("login")
	app.Post("/login", middleware.UnauthenticatedOnly, h.LoginSubmit)
	app.Get("/register", middleware.UnauthenticatedOnly, h.RegisterPage)
	app.Post("/register", middleware.UnauthenticatedOnly, h.RegisterSubmit)

	routes := app.Group("/auth")

	routes.Get("/", middleware.UnauthenticatedOnly, h.RedirectToLogin)
	routes.Get("/login", middleware.UnauthenticatedOnly, h.LoginPage)
	routes.Post("/login", middleware.UnauthenticatedOnly, h.LoginSubmit)
	routes.Get("/register", middleware.UnauthenticatedOnly, h.RegisterPage)
	routes.Post("/register", middleware.UnauthenticatedOnly, h.RegisterSubmit)

	// Authenticated only

	app.Get("/logout", middleware.AuthenticatedOnly, h.Logout).Name("logout")
	app.Post("/logout", middleware.AuthenticatedOnly, h.Logout)
	routes.Get("/logout", middleware.AuthenticatedOnly, h.Logout)
	routes.Post("/logout", middleware.AuthenticatedOnly, h.Logout)

	// TODO: Delete, Patch user
}

func (h *handler) RedirectToLogin(c *fiber.Ctx) error {
	return c.RedirectToRoute("login", nil)
}

func (h *handler) LoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"UsernameMinLength": model.MIN_USERNAME_LEN,
		"Config":            h.Config,
	}, "layouts/main")
}

func (h *handler) Logout(c *fiber.Ctx) error {

	common.ClearUser(c)
	if c.Get("HX-Request") != "" {
		return common.RedirectHTMX(c, "login", nil)
	}
	return c.RedirectToRoute("login", nil)
}

func (h *handler) LoginSubmit(c *fiber.Ctx) error {

	body := LoginRequestBody{}

	// parse body, unmarshall to AddTodoRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return c.Render("auth/forms/login", fiber.Map{
			"UsernameMinLength": model.MIN_USERNAME_LEN,
			"Errors":            []error{errors.New("invalid data")},
			"Config":            h.Config,
		})
	}

	user, err := model.LoginUser(h.DB, body.Username, body.Password)

	if err != nil {
		return c.Render("auth/forms/login", fiber.Map{
			"UsernameMinLength": model.MIN_USERNAME_LEN,
			"Errors":            []error{err},
			"Username":          body.Username,
			"Remember":          body.RememberMe,
			"Config":            h.Config,
		})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		common.USER_CLAIM_KEY: user.ID,
	}

	// define when to expire that token
	expires := time.Now().Add(time.Hour * 24 * 3) // expires in 3 days

	// this guy reaaaly wants us to remember him
	if body.RememberMe {
		// set expiry date to a looooooong time
		expires = time.Now().Add(time.Hour * 24 * 90) // expires in 90 days
	}

	// Create token
	claims["exp"] = expires.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	encryptedToken, err := token.SignedString([]byte(h.Config.EncryptionKey))
	if err != nil {
		return c.Render("auth/forms/login", fiber.Map{
			"UsernameMinLength": model.MIN_USERNAME_LEN,
			"Errors":            []error{errors.New("internal error")},
			"Username":          body.Username,
			"Remember":          body.RememberMe,
			"Config":            h.Config,
		})
	}

	// set cookie with encrypted token
	c.Cookie(&fiber.Cookie{
		Name:     common.USER_COOKIE_KEY,
		Value:    encryptedToken,
		SameSite: "lax",
		Expires:  expires,
		HTTPOnly: true,
	})

	// get all the todos
	todos, err := model.AllTodos(h.DB, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// redirect the guy to the main course
	c.Append("HX-Retarget", "#container")
	route, _ := c.GetRouteURL("todos", nil)
	c.Append("HX-Replace-Url", route)
	return c.Render("todos/index", fiber.Map{
		"Todos":    todos,
		"Username": user.Username,
		"Config":   h.Config,
	})
}
func (h *handler) RegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"UsernameMinLength": model.MIN_USERNAME_LEN,
		"Config":            h.Config,
	}, "layouts/main")
}

func (h *handler) RegisterSubmit(c *fiber.Ctx) error {
	body := RegisterRequestBody{}

	// parse body, unmarshall to AddTodoRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return c.Render("auth/forms/register", fiber.Map{
			"UsernameMinLength": model.MIN_USERNAME_LEN,
			"Errors":            []error{errors.New("invalid data")},
			"Config":            h.Config,
		})
	}

	_, err := model.CreateUser(h.DB, body.Username, body.Password, body.PasswordConfirm)

	if err != nil {
		return c.Render("auth/forms/register", fiber.Map{
			"UsernameMinLength": model.MIN_USERNAME_LEN,
			"Errors":            []error{err},
			"Username":          body.Username,
			"Config":            h.Config,
		})
	}

	// redirect the guy to the login page
	c.Append("HX-Retarget", "#container")
	return c.Render("auth/login", fiber.Map{
		"UsernameMinLength": model.MIN_USERNAME_LEN,
		"Username":          body.Username,
		"SuccessMessage":    "You can now login with your new account !",
		"Config":            h.Config,
	})
}
