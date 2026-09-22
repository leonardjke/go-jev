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

	out, err := json.Marshal(c)
	require.NoError(t, err)
	assert.NotContains(t, string(out), canaryKey)

	assert.Equal(t, canaryKey, string(c.APIKey))
}
