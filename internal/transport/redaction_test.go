package transport_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/leonardjke/go-jev/internal/privacy"
	"github.com/leonardjke/go-jev/internal/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const canaryKey = "sk-live-CANARY-d41d8cd98f00"

func TestClient_NeverPrintsAPIKey(t *testing.T) {
	c := transport.Client{
		Endpoint: "https://example.invalid",
		APIKey:   privacy.SensitiveString(canaryKey),
	}

	for _, format := range []string{"%v", "%+v", "%#v"} {
		assert.NotContains(t, fmt.Sprintf(format, c), canaryKey, format+" on Client")
		assert.NotContains(t, fmt.Sprintf(format, &c), canaryKey, format+" on *Client")
	}

	// Client is never marshaled in anger; this is a leak canary for APIKey losing
	// its type. HTTP is nil here, so the func fields SA1026 warns about never come up.
	out, err := json.Marshal(c) //nolint:gosec,musttag,staticcheck // see above
	require.NoError(t, err)
	assert.NotContains(t, string(out), canaryKey)

	assert.Equal(t, canaryKey, string(c.APIKey))
}
