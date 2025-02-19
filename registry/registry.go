package registry

import "github.com/gofiber/fiber/v2"

type RouteRegister func(app *fiber.App)

var modules []RouteRegister

func RegisterModule(register RouteRegister) { modules = append(modules, register) }

func RegisterAllModules(app *fiber.App) {
	for _, register := range modules {
		register(app)
	}
}
