package users

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ExecuteSeed performs the seeding process and returns a JSON response.
func ExecuteSeed(c *fiber.Ctx) error {
	if success := Seeders(); !success {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "seeder execution failed",
		})
	}
	return c.JSON(fiber.Map{
		"status": "seeder executed successfully",
	})
}

// CookieConfig1: Cookie compatible con peticiones cross-site
func CookieConfig1(c *fiber.Ctx) error {
	// Configurando una cookie persistente de 30 días con Fiber
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	c.Cookie(&fiber.Cookie{
		Name:     "test_cookie1",
		Value:    "value1",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		HTTPOnly: true,
		// Cambiamos a None para mayor compatibilidad con peticiones cross-site
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true, // Obligatorio con SameSite=None
	})

	// Forzar la respuesta con las cabeceras apropiadas
	c.Append("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Println("🍪 Configurando cookie1 con SameSite=None para compatibilidad cross-site")

	return c.JSON(fiber.Map{"message": "Cookie config 1 set", "cookie": "test_cookie1"})
}

// CookieConfig2: Prueba con SameSite None explícito
func CookieConfig2(c *fiber.Ctx) error {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	c.Cookie(&fiber.Cookie{
		Name:     "test_cookie2",
		Value:    "value2",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		HTTPOnly: true,
		// Cambiamos a None para resolver el problema de cross-site
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	c.Append("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Println("🍪 Configurando cookie2 con SameSite=None para compatibilidad cross-site")

	return c.JSON(fiber.Map{"message": "Cookie config 2 set", "cookie": "test_cookie2"})
}

// CookieConfig3: Cookie persistente con SameSite None
func CookieConfig3(c *fiber.Ctx) error {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	c.Cookie(&fiber.Cookie{
		Name:     "test_cookie3",
		Value:    "value3",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true, // Obligatorio para SameSite=None
	})

	c.Append("Cache-Control", "no-cache, no-store, must-revalidate")

	return c.JSON(fiber.Map{"message": "Cookie config 3 set", "cookie": "test_cookie3"})
}

// CookieConfig4: Modificado para usar también SameSite=None
func CookieConfig4(c *fiber.Ctx) error {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	c.Cookie(&fiber.Cookie{
		Name:     "test_cookie4",
		Value:    "value4",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		HTTPOnly: true,
		// Cambiamos a None para resolver el problema
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	c.Append("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Println("🍪 Configurando cookie4 con SameSite=None para compatibilidad cross-site")

	return c.JSON(fiber.Map{"message": "Cookie config 4 set", "cookie": "test_cookie4"})
}

// CookieConfig5: También modificado para usar SameSite=None
func CookieConfig5(c *fiber.Ctx) error {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	// Configuramos una cookie con una fecha de expiración muy lejana
	c.Cookie(&fiber.Cookie{
		Name:     "test_cookie5",
		Value:    "value5",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		HTTPOnly: true,
		// Cambiamos a None para resolver el problema
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	c.Append("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Println("🍪 Configurando cookie5 con SameSite=None para compatibilidad cross-site")

	return c.JSON(fiber.Map{"message": "Cookie config 5 set", "cookie": "test_cookie5"})
}

// CheckCookies: Endpoint para verificar las cookies recibidas del cliente con logging detallado
func CheckCookies(c *fiber.Ctx) error {
	// Obtener todas las cookies de la solicitud
	cookiesMap := make(map[string]string)

	// Obtener la cabecera completa para análisis
	cookieHeader := c.Get("Cookie")

	// Log detallado en consola para depuración
	fmt.Println("\n------ VERIFICACIÓN DE COOKIES ------")
	fmt.Printf("IP Remota: %s\n", c.IP())
	fmt.Printf("User-Agent: %s\n", c.Get("User-Agent"))
	fmt.Printf("Headers completos: %+v\n", c.GetReqHeaders())

	if cookieHeader == "" {
		fmt.Println("⚠️ NO SE RECIBIERON COOKIES")
		log.Println("Solicitud sin cookies recibida")
	} else {
		fmt.Printf("🍪 Cookie Header: %s\n", cookieHeader)
	}

	// Obtener cookies específicas para prueba
	testCookies := []string{
		"test_cookie1", "test_cookie2", "test_cookie3", "test_cookie4", "test_cookie5",
		"secure_cookie1", "secure_cookie2",
	}
	cookiesStatus := make(map[string]bool)

	fmt.Println("\n--- Verificando cookies individuales ---")

	// Verificar cada cookie de prueba e imprimir en consola
	for _, name := range testCookies {
		cookie := c.Cookies(name)
		cookiesMap[name] = cookie
		cookiesStatus[name] = cookie != ""

		if cookie != "" {
			fmt.Printf("✓ Cookie '%s' encontrada, valor: '%s'\n", name, cookie)
		} else {
			fmt.Printf("✗ Cookie '%s' NO encontrada\n", name)
		}
	}

	// Obtener otras posibles cookies
	fmt.Println("\n--- Otras cookies encontradas ---")
	allHeaders := c.GetReqHeaders()
	for key, values := range allHeaders {
		if key == "Cookie" {
			for _, value := range values {
				fmt.Printf("Header de Cookie: %s\n", value)
			}
		}
	}

	// Imprimimos un resumen
	fmt.Printf("\n--- RESUMEN ---\n")
	fmt.Printf("Total cookies esperadas: %d\n", len(testCookies))

	foundCount := 0
	for _, found := range cookiesStatus {
		if found {
			foundCount++
		}
	}

	fmt.Printf("Total cookies encontradas: %d\n", foundCount)
	fmt.Println("-----------------------------------")

	return c.JSON(fiber.Map{
		"message": "Cookies recibidas del cliente",
		"cookies": cookieHeader,
		"details": map[string]interface{}{
			"found":         foundCount > 0,
			"cookiesMap":    cookiesMap,
			"testStatus":    cookiesStatus,
			"ip":            c.IP(),
			"timestamp":     time.Now().Format(time.RFC3339),
			"totalExpected": len(testCookies),
			"totalFound":    foundCount,
			"allHeaders":    c.GetReqHeaders(),
		},
	})
}

// VerifyAuth comprueba si existen cookies de autenticación
func VerifyAuth(c *fiber.Ctx) error {
	// Verificamos si existe al menos una de nuestras cookies de prueba
	testCookies := []string{"test_cookie1", "test_cookie2", "test_cookie3", "test_cookie4", "test_cookie5"}
	authenticated := false
	foundCookies := []string{}

	for _, cookieName := range testCookies {
		cookie := c.Cookies(cookieName)
		if cookie != "" {
			authenticated = true
			foundCookies = append(foundCookies, cookieName)
		}
	}

	if authenticated {
		return c.JSON(fiber.Map{
			"authenticated": true,
			"message":       "Autenticación correcta",
			"cookies":       foundCookies,
		})
	}

	return c.Status(401).JSON(fiber.Map{
		"authenticated": false,
		"message":       "Cookies de prueba no encontradas",
	})
}

// SetSecureCookies configura múltiples cookies en una sola solicitud
func SetSecureCookies(c *fiber.Ctx) error {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	maxAge := int((30 * 24 * time.Hour).Seconds())

	// Modificamos para usar SameSite=None
	c.Cookie(&fiber.Cookie{
		Name:     "secure_cookie1",
		Value:    "test_value_1",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   maxAge,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "secure_cookie2",
		Value:    "test_value_2",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   maxAge,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	// Configurar cabeceras para evitar caché
	c.Append("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Append("Pragma", "no-cache")
	c.Append("Expires", "0")

	fmt.Println("🍪 Configurando secure_cookies con SameSite=None para compatibilidad cross-site")

	return c.JSON(fiber.Map{
		"message": "Múltiples cookies configuradas correctamente",
		"cookies": []string{"secure_cookie1", "secure_cookie2"},
	})
}

// NoPersistentCookie establece una cookie de sesión (sin MaxAge ni Expires)
func NoPersistentCookie(c *fiber.Ctx) error {
	// Cookie de sesión (sin MaxAge ni Expires)
	c.Cookie(&fiber.Cookie{
		Name:     "session_cookie",
		Value:    "session_value",
		Path:     "/",
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode, // Importante para cross-site
		Secure:   true,
	})

	fmt.Println("🍪 Configurando cookie de sesión con SameSite=None")

	return c.JSON(fiber.Map{
		"message": "Cookie de sesión configurada (no persistente)",
		"cookie":  "session_cookie",
	})
}

// PlainTextCookie establece una cookie sin seguridad para probar
func PlainTextCookie(c *fiber.Ctx) error {
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	// Cookie básica sin restricciones
	c.Cookie(&fiber.Cookie{
		Name:     "plain_cookie",
		Value:    "plain_value",
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		HTTPOnly: false, // Accesible desde JavaScript
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	fmt.Println("🍪 Configurando cookie simple sin HTTPOnly")

	return c.JSON(fiber.Map{
		"message": "Cookie sin HTTPOnly configurada",
		"cookie":  "plain_cookie",
	})
}

// RegisterRutasComplejas registra las rutas de prueba para las cookies.
func RegisterRutasComplejas(app *fiber.App) {
	group := app.Group("/api/complejas")
	group.Get("/config1", CookieConfig1)
	group.Get("/config2", CookieConfig2)
	group.Get("/config3", CookieConfig3)
	group.Get("/config4", CookieConfig4)
	group.Get("/config5", CookieConfig5)
	group.Get("/seed", ExecuteSeed)
	group.Get("/check-cookies", CheckCookies)
	group.Get("/verify-auth", VerifyAuth)
	group.Get("/set-secure-cookies", SetSecureCookies)

	// Nuevas rutas para probar otras configuraciones de cookies
	group.Get("/session-cookie", NoPersistentCookie)
	group.Get("/plain-cookie", PlainTextCookie)
}
