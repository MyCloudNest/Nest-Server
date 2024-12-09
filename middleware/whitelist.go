package middleware

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func WhitelistMiddleware(whitelistedIPs []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		for _, allowedIP := range whitelistedIPs {
			if isIPAllowed(clientIP, allowedIP) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"ok":      false,
			"message": "Access denied",
		})
	}
}

func isIPAllowed(clientIP, allowedIP string) bool {
	if strings.Contains(allowedIP, "/") {
		_, ipNet, err := net.ParseCIDR(allowedIP)
		if err != nil {
			return false
		}
		return ipNet.Contains(net.ParseIP(clientIP))
	}

	return clientIP == allowedIP
}
