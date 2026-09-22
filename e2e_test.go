//go:build e2e

package jev_test

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/leonardjke/go-jev"
	"github.com/stretchr/testify/assert"
)

func TestE2E_BasicQuestion(t *testing.T) {
	client := newIntegrationClient(t)

	ctx := t.Context()

	resp, err := client.Request(ctx,
		"Today in Warsaw the weather is a little rainy, in the end of the day there is a possibility of the sun",
		jev.Score("weather", "Is the rain a good weather", "good", "bad"),
	)
	assert.NoError(t, err)

	assert.NotEmpty(t, resp.Model)
	assert.NotEmpty(t, resp.Usage)
	assert.NotEmpty(t, resp.Answers["weather"])

	fmt.Println("resp: ", resp)
}

func newIntegrationClient(t *testing.T) *jev.Jev {
	t.Helper()

	apiKey := os.Getenv("API_KEY")
	apiUrl := os.Getenv("API_URL")

	if apiKey == "" {
		t.Skip("API_KEY was not set, skip integration test")
	}

	// The client only logs at debug level, and it logs to the test's own
	// output so the bodies show up under -v and stay attached to this test.
	logger := slog.New(slog.NewTextHandler(testWriter{t: t}, &slog.HandlerOptions{Level: slog.LevelDebug}))

	client, err := jev.New(apiKey, jev.WithEndpoint(apiUrl), jev.WithLogger(logger))
	assert.NoError(t, err)

	fmt.Println("client: ", client)

	return client
}

// testWriter routes what is written to it through t.Log, so log lines show up
// under -v and stay attached to the test that produced them. t.Output does
// this directly, but it landed in Go 1.25 and this module builds on 1.24.
type testWriter struct {
	t *testing.T
}

func (w testWriter) Write(p []byte) (int, error) {
	w.t.Helper()
	w.t.Log(strings.TrimSuffix(string(p), "\n"))
	return len(p), nil
}
