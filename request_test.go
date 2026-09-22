package jev_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leonardjke/go-jev"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, respBody string) (*jev.Jev, *[]byte) {
	t.Helper()

	var got []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		got = body

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, respBody)
	}))
	t.Cleanup(srv.Close)

	client, err := jev.New("test-key", jev.WithEndpoint(srv.URL))
	require.NoError(t, err)

	return client, &got
}

const stubResponse = `{"model":"jev-latest","answers":{"is_urgent":{"type":"noul","noul":0.9}},"usage":{"input_tokens":1,"output_tokens":2}}`

func TestRequest_WireShape(t *testing.T) {
	client, got := newTestClient(t, stubResponse)

	resp, err := client.Request(t.Context(), "the state",
		jev.Noul("is_urgent", "The message conveys urgency"),
		jev.Choice("department", "Which team should handle this", map[string]string{
			"billing":   "Payment or subscription issues",
			"technical": "Bugs or integration problems",
		}),
		jev.Score("amount_due_size", "How large is the `field` value?", "Small", "Typical", "Large").
			WithField(jev.Field{Name: "amount_due", Type: "number", Unit: "USD"}),
	)
	require.NoError(t, err)

	assert.JSONEq(t, `{
		"model": "jev-latest",
		"state": "the state",
		"questions": {
			"is_urgent": {
				"type": "noul",
				"instructions": "The message conveys urgency"
			},
			"department": {
				"type": "choice",
				"instructions": "Which team should handle this",
				"criteria": {
					"billing": "Payment or subscription issues",
					"technical": "Bugs or integration problems"
				}
			},
			"amount_due_size": {
				"type": "score",
				"instructions": {
					"question": "How large is the `+"`field`"+` value?",
					"field": {"name": "amount_due", "type": "number", "unit": "USD"}
				},
				"criteria": ["Small", "Typical", "Large"]
			}
		}
	}`, string(*got))

	require.NotNil(t, resp.Answers["is_urgent"].Noul)
	assert.Equal(t, 0.9, *resp.Answers["is_urgent"].Noul)
}

func TestRequest_NoulCriteriaAreOptional(t *testing.T) {
	client, got := newTestClient(t, stubResponse)

	_, err := client.Request(t.Context(), "s",
		jev.Noul("plain", "A plain check"),
		jev.Noul("described", "A described check").WithCriteria(map[string]string{
			"true":  "It happened",
			"false": "It did not",
		}),
	)
	require.NoError(t, err)

	var req map[string]any
	require.NoError(t, json.Unmarshal(*got, &req))

	questions := req["questions"].(map[string]any)
	assert.NotContains(t, questions["plain"], "criteria")
	assert.Contains(t, questions["described"], "criteria")
}

func TestRequest_Rejects(t *testing.T) {
	client, got := newTestClient(t, stubResponse)

	t.Run("no questions", func(t *testing.T) {
		_, err := client.Request(t.Context(), "s")
		assert.ErrorContains(t, err, "no questions to ask")
	})

	t.Run("empty key", func(t *testing.T) {
		_, err := client.Request(t.Context(), "s", jev.Noul("", "q"))
		assert.ErrorContains(t, err, "question key must not be empty")
	})

	t.Run("duplicate key", func(t *testing.T) {
		_, err := client.Request(t.Context(), "s",
			jev.Noul("same", "first"),
			jev.Score("same", "second", "a", "b"),
		)
		assert.ErrorContains(t, err, `duplicate question key "same"`)
	})

	t.Run("unencodable context value names the question", func(t *testing.T) {
		_, err := client.Request(t.Context(), "s",
			jev.Noul("ok", "fine"),
			jev.Noul("broken", "q").Set("tolerance", make(chan int)),
		)
		assert.ErrorContains(t, err, `question "broken"`)
		assert.ErrorContains(t, err, `key "tolerance"`)
	})

	assert.Empty(t, *got, "a rejected request must not reach the server")
}

func TestRequest_QuestionsDoNotAliasEachOther(t *testing.T) {
	client, got := newTestClient(t, stubResponse)

	base := jev.Noul("base", "Does `extracted_value` match?")
	withValue := base.WithExtractedValue("4471")

	_, err := client.Request(t.Context(), "s", base.WithField(jev.Field{Name: "a", Type: "string"}))
	require.NoError(t, err)
	assert.NotContains(t, string(*got), `"extracted_value":`, "branching off base must not leak the other branch's keys")

	_, err = client.Request(t.Context(), "s", withValue)
	require.NoError(t, err)
	assert.Contains(t, string(*got), `"extracted_value":`)
	assert.NotContains(t, string(*got), `"field":`)
}

func TestRequest_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"error":"slow down"}`)
	}))
	t.Cleanup(srv.Close)

	client, err := jev.New("test-key", jev.WithEndpoint(srv.URL))
	require.NoError(t, err)

	_, err = client.Request(t.Context(), "s", jev.Noul("q", "a question"))

	var apiErr *jev.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	assert.Contains(t, apiErr.Body, "slow down")
}
