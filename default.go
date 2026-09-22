package jev

import (
	"net/http"
	"time"

	"github.com/leonardjke/go-jev/internal/log"
)

const (
	defaultEndpoint = "https://api.typesafe.ai/v1/systemone"
	defaultModel    = "jev-latest"
	defaultTimeout  = 30 * time.Second

	maxResponseBytes = 8 << 20
)

var (
	defaultClient = &http.Client{Timeout: defaultTimeout}
	defaultLogger = &log.NoopLogger{}
)
