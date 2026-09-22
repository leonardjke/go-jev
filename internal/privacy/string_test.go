package privacy_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/leonardjke/go-jev/internal/privacy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

		j, err := json.Marshal(holder)
		require.NoError(t, err)
		assert.NotContains(t, string(j), secret)

		assert.NotContains(t, fmt.Sprintf("%+v", holder), secret)
	})

	assert.Equal(t, secret, string(ss))
}
