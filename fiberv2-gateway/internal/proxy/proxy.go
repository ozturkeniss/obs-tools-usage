package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// ProxyConfig holds configuration for the reverse proxy
type ProxyConfig struct {
	Timeout        time.Duration
	Retries        int
	RetryDelay     time.Duration
	StripPath      bool
	RewritePath    string
	AddHeaders     map[string]string
	RemoveHeaders  []string
}

// ReverseProxy handles reverse proxy functionality
type ReverseProxy struct {
	config ProxyConfig
	logger *logrus.Logger
}

// NewReverseProxy creates a new reverse proxy
func NewReverseProxy(config ProxyConfig, logger *logrus.Logger) *ReverseProxy {
	return &ReverseProxy{
		config: config,
		logger: logger,
	}
}

// ProxyRequest proxies a request to a backend server
func (rp *ReverseProxy) ProxyRequest(c *fiber.Ctx, backendURL string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: rp.config.Timeout,
	}

	// Build target URL
	targetURL := rp.buildTargetURL(backendURL, c.OriginalURL())

	// Create request
	req, err := rp.createRequest(c, targetURL)
	if err != nil {
		rp.logger.WithError(err).Error("Failed to create request")
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	// Execute request with retries
	var resp *http.Response
	for i := 0; i <= rp.config.Retries; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}

		if i < rp.config.Retries {
			rp.logger.WithFields(logrus.Fields{
				"attempt": i + 1,
				"url":     targetURL,
				"error":   err.Error(),
			}).Warn("Request failed, retrying")

			time.Sleep(rp.config.RetryDelay)
			continue
		}

		rp.logger.WithFields(logrus.Fields{
			"url":   targetURL,
			"error": err.Error(),
		}).Error("Request failed after all retries")

		return c.Status(502).JSON(fiber.Map{
			"error": "Backend service unavailable",
		})
	}

	defer resp.Body.Close()

	// Copy response headers
	rp.copyResponseHeaders(c, resp)

	// Copy response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rp.logger.WithError(err).Error("Failed to read response body")
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	// Set status code and body
	c.Status(resp.StatusCode)
	c.Set("Content-Type", resp.Header.Get("Content-Type"))

	return c.Send(body)
}

// buildTargetURL builds the target URL for the backend
func (rp *ReverseProxy) buildTargetURL(backendURL, originalURL string) string {
	targetURL := backendURL

	// Handle path rewriting
	if rp.config.StripPath && rp.config.RewritePath != "" {
		// Replace the original path with the rewrite path
		targetURL = strings.TrimSuffix(backendURL, "/") + rp.config.RewritePath
	} else if rp.config.StripPath {
		// Strip the path from the original URL
		targetURL = backendURL
	} else {
		// Use the original URL path
		if !strings.HasSuffix(backendURL, "/") && !strings.HasPrefix(originalURL, "/") {
			targetURL = backendURL + "/" + originalURL
		} else {
			targetURL = backendURL + originalURL
		}
	}

	return targetURL
}

// createRequest creates an HTTP request from the Fiber context
func (rp *ReverseProxy) createRequest(c *fiber.Ctx, targetURL string) (*http.Request, error) {
	// Create request body reader
	var bodyReader io.Reader
	if len(c.Body()) > 0 {
		bodyReader = bytes.NewReader(c.Body())
	}

	// Create request
	req, err := http.NewRequestWithContext(c.Context(), c.Method(), targetURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Copy headers from original request
	rp.copyRequestHeaders(c, req)

	// Add custom headers
	for key, value := range rp.config.AddHeaders {
		req.Header.Set(key, value)
	}

	// Remove specified headers
	for _, header := range rp.config.RemoveHeaders {
		req.Header.Del(header)
	}

	return req, nil
}

// copyRequestHeaders copies headers from Fiber context to HTTP request
func (rp *ReverseProxy) copyRequestHeaders(c *fiber.Ctx, req *http.Request) {
	// Copy all headers from the original request
	c.Request().Header.VisitAll(func(key, value []byte) {
		headerName := string(key)
		headerValue := string(value)

		// Skip headers that shouldn't be forwarded
		if rp.shouldSkipHeader(headerName) {
			return
		}

		req.Header.Set(headerName, headerValue)
	})

	// Set X-Forwarded-For header
	if clientIP := c.IP(); clientIP != "" {
		req.Header.Set("X-Forwarded-For", clientIP)
	}

	// Set X-Forwarded-Proto header
	req.Header.Set("X-Forwarded-Proto", c.Protocol())

	// Set X-Forwarded-Host header
	req.Header.Set("X-Forwarded-Host", c.Hostname())
}

// copyResponseHeaders copies headers from HTTP response to Fiber context
func (rp *ReverseProxy) copyResponseHeaders(c *fiber.Ctx, resp *http.Response) {
	for key, values := range resp.Header {
		// Skip headers that shouldn't be copied
		if rp.shouldSkipResponseHeader(key) {
			continue
		}

		// Set header with all values
		for _, value := range values {
			c.Append(key, value)
		}
	}
}

// shouldSkipHeader determines if a header should be skipped when forwarding
func (rp *ReverseProxy) shouldSkipHeader(headerName string) bool {
	skipHeaders := []string{
		"Connection",
		"Upgrade",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
	}

	headerName = strings.ToLower(headerName)
	for _, skip := range skipHeaders {
		if headerName == strings.ToLower(skip) {
			return true
		}
	}

	return false
}

// shouldSkipResponseHeader determines if a response header should be skipped
func (rp *ReverseProxy) shouldSkipResponseHeader(headerName string) bool {
	skipHeaders := []string{
		"Connection",
		"Upgrade",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
	}

	headerName = strings.ToLower(headerName)
	for _, skip := range skipHeaders {
		if headerName == strings.ToLower(skip) {
			return true
		}
	}

	return false
}

// ProxyMiddleware creates a middleware function for proxying requests
func (rp *ReverseProxy) ProxyMiddleware(backendURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		
		// Log request
		rp.logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"url":    c.OriginalURL(),
			"backend": backendURL,
		}).Info("Proxying request")

		// Proxy the request
		err := rp.ProxyRequest(c, backendURL)
		
		// Log response
		duration := time.Since(startTime)
		rp.logger.WithFields(logrus.Fields{
			"method":   c.Method(),
			"path":     c.Path(),
			"status":   c.Response().StatusCode(),
			"duration": duration,
			"backend":  backendURL,
		}).Info("Request proxied")

		return err
	}
}

// FastHTTPProxy proxies using FastHTTP for better performance
func (rp *ReverseProxy) FastHTTPProxy(c *fiber.Ctx, backendURL string) error {
	// Create FastHTTP client
	client := &fasthttp.Client{
		ReadTimeout:  rp.config.Timeout,
		WriteTimeout: rp.config.Timeout,
	}

	// Create request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Set request URL
	req.SetRequestURI(rp.buildTargetURL(backendURL, c.OriginalURL()))
	req.Header.SetMethod(c.Method())

	// Copy request headers
	c.Request().Header.VisitAll(func(key, value []byte) {
		if !rp.shouldSkipHeader(string(key)) {
			req.Header.Set(string(key), string(value))
		}
	})

	// Set request body
	if len(c.Body()) > 0 {
		req.SetBody(c.Body())
	}

	// Add custom headers
	for key, value := range rp.config.AddHeaders {
		req.Header.Set(key, value)
	}

	// Remove specified headers
	for _, header := range rp.config.RemoveHeaders {
		req.Header.Del(header)
	}

	// Execute request
	err := client.Do(req, resp)
	if err != nil {
		rp.logger.WithFields(logrus.Fields{
			"url":   backendURL,
			"error": err.Error(),
		}).Error("FastHTTP request failed")

		return c.Status(502).JSON(fiber.Map{
			"error": "Backend service unavailable",
		})
	}

	// Copy response headers
	resp.Header.VisitAll(func(key, value []byte) {
		if !rp.shouldSkipResponseHeader(string(key)) {
			c.Set(string(key), string(value))
		}
	})

	// Set status and body
	c.Status(resp.StatusCode())
	c.Set("Content-Type", string(resp.Header.Peek("Content-Type")))

	return c.Send(resp.Body())
}
