package network

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	MaxBackoff time.Duration

	MaxRetries int

	Timeout time.Duration

	client *http.Client
}
