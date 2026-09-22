package data

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Instructions is what a Question asks. On the wire it is either a bare
// string or an object with a required "question" key plus any number of
// context keys ("field", "extracted_value", "compare", "focus", "note", ...).
//
// Question is the one key that is always present, so it is typed. Everything
// else is kept as raw JSON in Extra: the context keys are an open set, and
// keys this library doesn't know about must survive a decode/encode round
// trip untouched.
type Instructions struct {
	Question string
	Extra    map[string]json.RawMessage

	// err records the first failure from a With* call so the fluent chain
	// stays expression-shaped; MarshalJSON reports it.
	err error
}

// Field describes the value a question is asked about. It is stored under the
// "field" key rather than being a member of Instructions, because not every
// question is about a field.
type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Unit        string `json:"unit,omitempty"`
	Description string `json:"description,omitempty"`
}

// Ask starts Instructions from the question text. Used alone it encodes as a
// bare JSON string; each With* call adds a context key and switches the
// encoding to an object.
func Ask(question string) Instructions {
	return Instructions{Question: question}
}

// Set stores value under an arbitrary key. It is the escape hatch for context
// keys this library has no typed helper for yet.
func (i Instructions) Set(key string, value any) Instructions {
	raw, err := json.Marshal(value)
	if err != nil {
		if i.err == nil {
			i.err = fmt.Errorf("instructions: key %q: %w", key, err)
		}
		return i
	}

	// Copy so chains branched off a shared Instructions don't alias.
	extra := make(map[string]json.RawMessage, len(i.Extra)+1)
	maps.Copy(extra, i.Extra)
	extra[key] = raw
	i.Extra = extra

	return i
}

// WithField sets "field": the value being judged.
func (i Instructions) WithField(f Field) Instructions {
	return i.Set("field", f)
}

// WithExtractedValue sets "extracted_value": the candidate value to check the
// field against.
func (i Instructions) WithExtractedValue(value any) Instructions {
	return i.Set("extracted_value", value)
}

// WithCompare sets "compare": the state paths the question weighs against
// each other, e.g. "ticket.sender.display_name".
func (i Instructions) WithCompare(paths ...string) Instructions {
	return i.Set("compare", paths)
}

// WithFocus sets "focus": what the judge should look at.
func (i Instructions) WithFocus(focus string) Instructions {
	return i.Set("focus", focus)
}

// WithNote sets "note": a caveat for how to judge.
func (i Instructions) WithNote(note string) Instructions {
	return i.Set("note", note)
}

// Err reports the first failure from a Set or With* call on these
// Instructions. The fluent chain stays expression-shaped, so the error is
// carried here rather than returned; callers assembling a request report it
// alongside the key of the question it belongs to.
func (i Instructions) Err() error {
	return i.err
}

// Get returns the raw JSON stored under key.
func (i Instructions) Get(key string) (json.RawMessage, bool) {
	raw, ok := i.Extra[key]
	return raw, ok
}

// Decode unmarshals the value stored under key into v, e.g.
// var f Field; err := ins.Decode("field", &f).
func (i Instructions) Decode(key string, v any) error {
	raw, ok := i.Extra[key]
	if !ok {
		return fmt.Errorf("instructions: no key %q", key)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("instructions: key %q: %w", key, err)
	}
	return nil
}

// MarshalJSON encodes a bare string when there is no context, an object
// otherwise, so simple questions stay simple on the wire.
func (i Instructions) MarshalJSON() ([]byte, error) {
	if i.err != nil {
		return nil, i.err
	}
	if len(i.Extra) == 0 {
		return json.Marshal(i.Question)
	}

	out := make(map[string]json.RawMessage, len(i.Extra)+1)
	maps.Copy(out, i.Extra)

	question, err := json.Marshal(i.Question)
	if err != nil {
		return nil, err
	}
	out["question"] = question

	return json.Marshal(out)
}

// UnmarshalJSON accepts either form, keeping unrecognized context keys in Extra.
func (i *Instructions) UnmarshalJSON(data []byte) error {
	var question string
	if err := json.Unmarshal(data, &question); err == nil {
		i.Question = question
		i.Extra = nil
		return nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("instructions: want string or object: %w", err)
	}

	if q, ok := raw["question"]; ok {
		if err := json.Unmarshal(q, &i.Question); err != nil {
			return fmt.Errorf("instructions: question: %w", err)
		}
		delete(raw, "question")
	}

	if len(raw) == 0 {
		raw = nil
	}
	i.Extra = raw

	return nil
}
