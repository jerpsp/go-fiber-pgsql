package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// isJSON checks if a string is a valid JSON
func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// LoggerConfig holds the configuration for the logger middleware
type LoggerConfig struct {
	// Skip sensitive routes from logging bodies
	SkipSensitiveRoutes []string
}

// DefaultLoggerConfig returns the default logger configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		SkipSensitiveRoutes: []string{"/api/v1/auth/signin", "/api/v1/auth/signup", "/api/v1/auth/refresh"},
	}
}

// Logger returns a middleware function that logs information about each HTTP request and response.
func Logger() fiber.Handler {
	// Use default config
	config := DefaultLoggerConfig()
	return LoggerWithConfig(config)
}

// LoggerWithConfig returns a logger middleware with custom configuration
func LoggerWithConfig(config LoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Get the path
		path := c.Path()

		// Get the HTTP method
		method := c.Method()

		// Get request body (store a copy before processing)
		reqBody := string(c.Body())

		// Store the original body since Fiber doesn't have built-in body cloning
		// We need to manually store the body before passing to the next handler
		var reqBodyMap map[string]interface{}

		// Check content type to properly handle different types of request data
		contentType := string(c.Request().Header.ContentType())
		isFormData := strings.Contains(contentType, "multipart/form-data") ||
			strings.Contains(contentType, "application/x-www-form-urlencoded")

		// Only try to unmarshal JSON if content type is JSON and body is valid JSON
		if reqBody != "" && !isFormData && isJSON(reqBody) {
			json.Unmarshal([]byte(reqBody), &reqBodyMap)
		}

		// Process request
		err := c.Next()

		// Calculate request processing time
		duration := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Get IP
		ip := c.IP()

		// Get content length of response
		contentLength := len(c.Response().Body())

		// Format and color status code
		var statusColor, methodColor, resetColor string

		// ANSI color codes
		resetColor = "\033[0m"

		// Method color
		switch method {
		case "GET":
			methodColor = "\033[32m" // Green
		case "POST":
			methodColor = "\033[34m" // Blue
		case "PUT":
			methodColor = "\033[33m" // Yellow
		case "DELETE":
			methodColor = "\033[31m" // Red
		case "PATCH":
			methodColor = "\033[36m" // Cyan
		default:
			methodColor = "\033[0m" // Reset
		}

		// Status color
		switch {
		case status >= 200 && status < 300:
			statusColor = "\033[32m" // Green
		case status >= 300 && status < 400:
			statusColor = "\033[33m" // Yellow
		case status >= 400 && status < 500:
			statusColor = "\033[31m" // Red
		default:
			statusColor = "\033[31m" // Red
		}

		// loc, _ := time.LoadLocation("Asia/Bangkok")
		now := time.Now().Format("2006-01-02 15:04:05")

		// Log format for request/response summary
		logMessage := fmt.Sprintf(
			"%s | %s[%s]%s%s | %s |%s %s%d%s | %s | %dB | %s",
			now,
			methodColor, method, resetColor,
			resetColor, path, resetColor,
			statusColor, status, resetColor,
			duration.String(),
			contentLength,
			ip,
		)

		// Log the request/response summary
		fmt.Println(logMessage)

		// Check if this is a sensitive route
		isSensitiveRoute := false
		currentPath := c.Path()
		for _, route := range config.SkipSensitiveRoutes {
			if strings.HasPrefix(currentPath, route) {
				isSensitiveRoute = true
				break
			}
		}

		// Get and log request body if it exists, isn't empty, and isn't a sensitive route
		if reqBody != "" && !isSensitiveRoute {
			// Handle different content types appropriately
			contentType := string(c.Request().Header.ContentType())

			if strings.Contains(contentType, "multipart/form-data") {
				// For multipart form data (file uploads), don't print the raw binary data
				fmt.Println("Request Body: [MULTIPART FORM DATA / FILE UPLOAD]")

				// Optionally log form field names without values
				form, err := c.MultipartForm()
				if err == nil {
					fields := make([]string, 0, len(form.Value))
					for field := range form.Value {
						fields = append(fields, field)
					}

					files := make([]string, 0, len(form.File))
					for file := range form.File {
						files = append(files, file)
					}

					if len(fields) > 0 {
						fmt.Printf("  Form Fields: %s\n", strings.Join(fields, ", "))
					}
					if len(files) > 0 {
						fmt.Printf("  File Fields: %s\n", strings.Join(files, ", "))
					}
				}
			} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				// For URL-encoded form data, print field names and values
				fmt.Println("Request Body: [URL-ENCODED FORM DATA]")

				// Parse form data
				form := c.Context().QueryArgs()
				form.VisitAll(func(key, value []byte) {
					fmt.Printf("  %s: %s\n", string(key), string(value))
				})
			} else {
				// For JSON and other content types, use the formatJSON function
				fmt.Printf("Request Body: %s\n", formatJSON(reqBody))
			}
		} else if reqBody != "" && isSensitiveRoute {
			fmt.Println("Request Body: [REDACTED - SENSITIVE DATA]")
		}

		// Get and log response body if it exists, isn't empty, and isn't a sensitive route
		// resBody := string(c.Response().Body())
		// if resBody != "" && !isSensitiveRoute {
		// 	fmt.Printf("Response Body: %s\n", formatJSON(resBody))
		// } else if resBody != "" && isSensitiveRoute {
		// 	fmt.Println("Response Body: [REDACTED - SENSITIVE DATA]")
		// }

		// Add separator for better readability between requests
		fmt.Println(strings.Repeat("-", 100))

		return err
	}
}

// formatJSON formats JSON string for better readability
func formatJSON(jsonStr string) string {
	// Check if it's valid JSON
	if !isJSON(jsonStr) {
		return jsonStr
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(jsonStr), "", "  "); err != nil {
		return jsonStr // Return original if formatting fails
	}

	return prettyJSON.String()
}
