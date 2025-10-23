module fiberv2-gateway

go 1.22

require (
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/gofiber/contrib/otelfiber v1.0.10
	github.com/gofiber/contrib/prometheus v1.1.0
	github.com/sony/gobreaker v0.5.0
	github.com/valyala/fasthttp v1.53.0
	github.com/prometheus/client_golang v1.19.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.9.0
	go.uber.org/ratelimit v0.3.1
	golang.org/x/time v0.5.0
)