package data

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Model     string              `json:"model"`
	State     string              `json:"state"`
	Questions map[string]Question `json:"questions"`
}

type QuestionType string

const (
	QuestionNoul   QuestionType = "noul"
	QuestionChoice QuestionType = "choice"
	QuestionScore  QuestionType = "score"
)

type Question struct {
	Type         QuestionType `json:"type"`
	Instructions Instructions `json:"instructions"`
	Criteria     any          `json:"criteria,omitempty"`
}

// CriteriaScore is the criteria shape for "score" questions: an ordered list
// of labels from lowest to highest, e.g. ["Calm", "Frustrated", "Very angry"].
type CriteriaScore []string

// CriteriaOthers is the criteria shape for "noul" and "choice" questions: a
// set of possible answers mapped to the description that qualifies them.
type CriteriaOthers map[string]string

// UnmarshalJSON decodes Criteria into the concrete type that matches Type
// (CriteriaScore for "score", CriteriaOthers for "noul"/"choice"), so callers
// get a typed value back from QuestionCriteria instead of asserting on `any`.
func (q *Question) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type         QuestionType    `json:"type"`
		Instructions Instructions    `json:"instructions"`
		Criteria     json.RawMessage `json:"criteria"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	q.Type = raw.Type
	q.Instructions = raw.Instructions

	switch raw.Type {
	case QuestionScore:
		var c CriteriaScore
		if err := json.Unmarshal(raw.Criteria, &c); err != nil {
			return fmt.Errorf("question type %q: criteria: %w", raw.Type, err)
		}
		q.Criteria = c
	case QuestionNoul:
		// "noul" questions don't require criteria (e.g. a plain yes/no check).
		if len(raw.Criteria) == 0 {
			break
		}
		var c CriteriaOthers
		if err := json.Unmarshal(raw.Criteria, &c); err != nil {
			return fmt.Errorf("question type %q: criteria: %w", raw.Type, err)
		}
		q.Criteria = c
	case QuestionChoice:
		var c CriteriaOthers
		if err := json.Unmarshal(raw.Criteria, &c); err != nil {
			return fmt.Errorf("question type %q: criteria: %w", raw.Type, err)
		}
		q.Criteria = c
	default:
		return fmt.Errorf("question: unknown type %q", raw.Type)
	}

	return nil
}

// QuestionCriteria type-asserts a Question's criteria into the wanted
// concrete type. Use QuestionCriteria[CriteriaScore](q) for a "score"
// question or QuestionCriteria[CriteriaOthers](q) for a "noul"/"choice" one;
// ok is false if q.Criteria isn't that type (e.g. wrong T for q.Type).
func QuestionCriteria[T CriteriaScore | CriteriaOthers](q Question) (T, bool) {
	c, ok := q.Criteria.(T)
	return c, ok
}
