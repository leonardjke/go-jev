package jev

import (
	"fmt"

	"github.com/leonardjke/go-jev/data"
)

// Field describes the value a question is asked about.
type Field = data.Field

// Question is one question to ask, along with the key its answer comes back
// under in Response.Answers.
//
// Build one with Noul, Choice or Score. Those cover the common case on their
// own; the With* methods add context to the questions that need it, and leave
// the ones that don't alone:
//
//	jev.Noul("is_urgent", "The message conveys urgency")
//
//	jev.Score("amount_due_size", "How large is the `field` value in `source_text`?",
//		"Small", "Typical", "Unusually large").
//		WithField(jev.Field{Name: "amount_due", Type: "number", Unit: "USD"})
//
// The zero Question is not usable; Request rejects one with an empty key.
type Question struct {
	key string
	q   data.Question
}

// Noul asks a yes/no question. Criteria are optional — a plain check needs
// none, WithCriteria describes what makes each answer true or false.
func Noul(key, question string) Question {
	return Question{
		key: key,
		q: data.Question{
			Type:         data.QuestionNoul,
			Instructions: data.Ask(question),
		},
	}
}

// Choice asks the judge to pick one key out of criteria, which maps each
// possible answer to the description that qualifies it.
func Choice(key, question string, criteria map[string]string) Question {
	return Question{
		key: key,
		q: data.Question{
			Type:         data.QuestionChoice,
			Instructions: data.Ask(question),
			Criteria:     data.CriteriaOthers(criteria),
		},
	}
}

// Score asks the judge to pick one point on an ordered scale. The scale
// labels run lowest to highest, e.g. "Calm", "Frustrated", "Very angry".
func Score(key, question string, scale ...string) Question {
	return Question{
		key: key,
		q: data.Question{
			Type:         data.QuestionScore,
			Instructions: data.Ask(question),
			Criteria:     data.CriteriaScore(scale),
		},
	}
}

// Key returns the key this question's answer comes back under.
func (q Question) Key() string {
	return q.key
}

// WithCriteria describes what makes each answer of a Noul question true or
// false. It replaces any criteria already set.
func (q Question) WithCriteria(criteria map[string]string) Question {
	q.q.Criteria = data.CriteriaOthers(criteria)
	return q
}

// WithField sets "field": the value being judged.
func (q Question) WithField(f Field) Question {
	return q.set(func(i data.Instructions) data.Instructions { return i.WithField(f) })
}

// WithExtractedValue sets "extracted_value": the candidate value to check the
// field against.
func (q Question) WithExtractedValue(value any) Question {
	return q.set(func(i data.Instructions) data.Instructions { return i.WithExtractedValue(value) })
}

// WithCompare sets "compare": the state paths the question weighs against
// each other, e.g. "ticket.sender.display_name".
func (q Question) WithCompare(paths ...string) Question {
	return q.set(func(i data.Instructions) data.Instructions { return i.WithCompare(paths...) })
}

// WithFocus sets "focus": what the judge should look at.
func (q Question) WithFocus(focus string) Question {
	return q.set(func(i data.Instructions) data.Instructions { return i.WithFocus(focus) })
}

// WithNote sets "note": a caveat for how to judge.
func (q Question) WithNote(note string) Question {
	return q.set(func(i data.Instructions) data.Instructions { return i.WithNote(note) })
}

// Set stores value under an arbitrary instruction key. It is the escape hatch
// for context keys this library has no typed helper for yet.
func (q Question) Set(key string, value any) Question {
	return q.set(func(i data.Instructions) data.Instructions { return i.Set(key, value) })
}

func (q Question) Validate() error {
	if err := q.q.Instructions.Err(); err != nil {
		return fmt.Errorf("question %q: %w", q.key, err)
	}
	return nil
}

func (q Question) set(fn func(data.Instructions) data.Instructions) Question {
	q.q.Instructions = fn(q.q.Instructions)
	return q
}
