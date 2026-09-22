package jev

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/leonardjke/go-jev/internal/privacy"
	"github.com/leonardjke/go-jev/internal/transport"
)

type Jev struct { //nolint:recvcheck // Format takes a value receiver on purpose; see its doc comment.
	transport *transport.Client

	model string
}

func New(apiKey string, opts ...Option) (*Jev, error) {
	apiKey = strings.TrimSpace(apiKey)

	j := Jev{
		transport: &transport.Client{
			HTTP:             defaultClient,
			Logger:           defaultLogger,
			Endpoint:         defaultEndpoint,
			MaxResponseBytes: maxResponseBytes,
			APIKey:           privacy.SensitiveString(apiKey),
		},
		model: defaultModel,
	}

	for _, opt := range opts {
		opt(&j)
	}

	if j.transport.APIKey == "" {
		return nil, errors.New("jev api key is required")
	}

	if j.transport.Logger == nil {
		return nil, errors.New("jev logger should not be nil")
	}

	return &j, nil
}

// Format keeps the client printable without printing the credential. The
// receiver is a value so that a copied Jev is covered too, rather than
// falling back to reflection over its fields.
func (j Jev) Format(f fmt.State, _ rune) {
	if j.transport == nil {
		_, _ = io.WriteString(f, "Client{uninitialized}")
		return
	}

	_, _ = fmt.Fprintf(f, "Client{endpoint: %q, model: %q}", j.transport.Endpoint, j.model)
}
