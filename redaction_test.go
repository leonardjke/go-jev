package jev_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leonardjke/go-jev"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const canaryKey = "sk-live-CANARY-d41d8cd98f00"

func TestClient_NeverPrintsAPIKey(t *testing.T) {
	client, err := jev.New(canaryKey, jev.WithEndpoint("https://example.invalid"))
	require.NoError(t, err)

	t.Run("format verbs", func(t *testing.T) {
		for _, format := range []string{"%v", "%s", "%+v", "%#v", "%d"} {
			assert.NotContains(t, fmt.Sprintf(format, client), canaryKey, format+" on *Jev")
			assert.NotContains(t, fmt.Sprintf(format, *client), canaryKey, format+" on Jev")
		}

		assert.Contains(t, fmt.Sprintf("%v", client), "endpoint", "the client should still print something useful")
	})

	t.Run("marshaled", func(t *testing.T) {
		// Jev has no exported fields today; the canary is that it stays that way.
		out, err := json.Marshal(client) //nolint:staticcheck // SA9005: deliberate.
		require.NoError(t, err)
		assert.NotContains(t, string(out), canaryKey)
	})

	t.Run("uninitialized", func(t *testing.T) {
		var zero jev.Jev
		assert.Equal(t, "Client{uninitialized}", fmt.Sprintf("%v", zero))
	})
}

func TestErrors_NeverCarryAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"error":"slow down"}`)
	}))
	t.Cleanup(srv.Close)

	cases := map[string]string{
		"api error":    srv.URL,
		"dial failure": "http://127.0.0.1:1/",
		"bad url":      "://bad",
	}

	for name, endpoint := range cases {
		t.Run(name, func(t *testing.T) {
			client, err := jev.New(canaryKey, jev.WithEndpoint(endpoint))
			require.NoError(t, err)

			_, err = client.Request(t.Context(), "s", jev.Noul("q", "a question"))
			require.Error(t, err)
			assert.NotContains(t, err.Error(), canaryKey)
			assert.NotContains(t, fmt.Sprintf("%+v", err), canaryKey)
		})
	}
}

func TestDebugLog_NeverCarriesAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, stubResponse)
	}))
	t.Cleanup(srv.Close)

	var logged bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logged, &slog.HandlerOptions{Level: slog.LevelDebug}))

	client, err := jev.New(canaryKey, jev.WithEndpoint(srv.URL), jev.WithLogger(logger))
	require.NoError(t, err)

	_, err = client.Request(t.Context(), "s", jev.Noul("is_urgent", "The message conveys urgency"))
	require.NoError(t, err)

	// The client logged through us at all — otherwise this test proves nothing.
	require.NotEmpty(t, logged.String(), "expected the client to log at debug level")
	assert.NotContains(t, logged.String(), canaryKey)

	logger.Debug("client", "client", client)
	assert.NotContains(t, logged.String(), canaryKey, "logging the client itself must stay clean")
}

func TestRequestBody_NeverCarriesAPIKey(t *testing.T) {
	var body, authorization string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		body, authorization = string(raw), r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, stubResponse)
	}))
	t.Cleanup(srv.Close)

	client, err := jev.New(canaryKey, jev.WithEndpoint(srv.URL))
	require.NoError(t, err)

	_, err = client.Request(t.Context(), "s", jev.Noul("is_urgent", "The message conveys urgency"))
	require.NoError(t, err)

	assert.NotContains(t, body, canaryKey)
	assert.Equal(t, "Bearer "+canaryKey, authorization)
	assert.True(t, strings.HasPrefix(authorization, "Bearer "))
}
