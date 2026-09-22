package jev

import (
	"context"
	"errors"
	"fmt"

	"github.com/leonardjke/go-jev/data"
	"github.com/leonardjke/go-jev/internal/transport"
)

// Request judges state against questions, using the model the client was
// built with.
//
//	resp, err := client.Request(ctx, ticket,
//		jev.Choice("department", "Which team should handle this?", map[string]string{
//			"billing":   "Payment or subscription issues",
//			"technical": "Bugs or integration problems",
//		}),
//		jev.Score("frustration", "How frustrated is the customer?",
//			"Calm", "Frustrated but civil", "Very angry"),
//	)
//
// Every question needs a key of its own; the answers come back under those
// keys in Response.Answers.
func (j *Jev) Request(ctx context.Context, state string, questions ...Question) (*data.Response, error) {
	if len(questions) == 0 {
		return nil, errors.New("jev: no questions to ask")
	}

	asked := make(map[string]data.Question, len(questions))
	for _, q := range questions {
		if q.key == "" {
			return nil, errors.New("jev: question key must not be empty")
		}
		if _, dup := asked[q.key]; dup {
			return nil, fmt.Errorf("jev: duplicate question key %q", q.key)
		}
		if err := q.Validate(); err != nil {
			return nil, fmt.Errorf("jev: %w", err)
		}
		asked[q.key] = q.q
	}

	return j.do(ctx, data.Request{
		Model:     j.model,
		State:     state,
		Questions: asked,
	})
}

func (j *Jev) do(ctx context.Context, req data.Request) (*data.Response, error) {
	var resp data.Response

	if err := j.transport.Post(ctx, req, &resp); err != nil {
		var status *transport.StatusError
		if errors.As(err, &status) {
			return nil, &APIError{StatusCode: status.StatusCode, Body: status.Body}
		}

		return nil, err
	}

	return &resp, nil
}
