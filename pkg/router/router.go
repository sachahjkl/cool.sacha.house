package router

import (
	"fmt"
	"sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/config"
	"sachahjkl/htmx_go/pkg/middleware"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"

	"github.com/gofiber/template/html/v2"
)

func New(c *config.Config) *fiber.App {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		ServerHeader: "Fiber",
		AppName:      "HTMX + Go",
		Views:        engine,
		// Views Layout is the global layout for all template render until override on Render function.
		// ViewsLayout: "layouts/main",
	})

	// Some cool metrics
	app.Get("/metrics", monitor.New())

	// Faster favicon
	app.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
		URL:  "/assets/favicon.ico",
	}))

	// Log in console
	app.Use(logger.New())

	// Compress request
	app.Use(compress.New())

	// Protect server
	// app.Use(helmet.New())

	// use jwt middleware for authentication
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(c.EncryptionKey)},
		TokenLookup: fmt.Sprintf("cookie:%v", common.USER_COOKIE_KEY), // should be "cookie:user"
		ContextKey:  common.USER_CONTEXT_KEY,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Next() // let me do the checking in my own middleware
		},
	}))

	return app
}

func RegisterDefaultRoutes(app *fiber.App, db *gorm.DB) {
	app.Static("/assets", "./assets", fiber.Static{
		Browse: true,
	})
	m := middleware.New(db)
	app.Get("/", m.WithUser, func(c *fiber.Ctx) error {
		user := c.Locals(common.USER_CONTEXT_KEY)

		if user != nil {
			return c.RedirectToRoute("todos", nil)
		} else {
			return c.RedirectToRoute("login", nil)
		}

	}).Name("root")
}
