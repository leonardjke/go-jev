package privacy_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/leonardjke/go-jev/internal/privacy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gosec // G101: a stand-in key for the test, not a real credential.
const secret = "sk-live-supersecret"

func TestSensitiveString_Redacts(t *testing.T) {
	ss := privacy.SensitiveString(secret)

	t.Run("verbs", func(t *testing.T) {
		for _, format := range []string{"%v", "%s", "%q", "%#v", "%d"} {
			out := fmt.Sprintf(format, ss)
			assert.Equal(t, "[REDACTED]", out, format)
		}
	})

	t.Run("marshalers", func(t *testing.T) {
		j, err := json.Marshal(ss)
		require.NoError(t, err)
		assert.JSONEq(t, `"[REDACTED]"`, string(j))

		text, err := ss.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, "[REDACTED]", string(text))
	})

	t.Run("as a struct field", func(t *testing.T) {
		holder := struct {
			APIKey privacy.SensitiveString
		}{ss}

		j, err := json.Marshal(holder) //nolint:gosec,musttag // marshaling the secret-shaped field is the point of the test.
		require.NoError(t, err)
		assert.NotContains(t, string(j), secret)

		assert.NotContains(t, fmt.Sprintf("%+v", holder), secret)
	})

	assert.Equal(t, secret, string(ss))
}
