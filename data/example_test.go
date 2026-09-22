package data_test

import (
	"encoding/json"
	"fmt"

	"github.com/leonardjke/go-jev/data"
)

// ExampleRequest covers the four instruction shapes at the wire level: plain
// text, a field plus the value extracted for it, a field on its own, and a
// question with extra guidance for the judge. Callers assembling a request
// normally use the constructors in the jev package instead; this is the layer
// underneath them.
func ExampleRequest() {
	req := data.Request{
		Model: "jev-latest",
		State: "...",
		Questions: map[string]data.Question{
			// Plain text — encodes as a bare JSON string.
			"department": {
				Type:         data.QuestionChoice,
				Instructions: data.Ask("Which team should handle this"),
				Criteria: data.CriteriaOthers{
					"billing":   "Payment or subscription issues",
					"technical": "Bugs or integration problems",
				},
			},

			// Field + the value we extracted for it.
			"invoice_number_matches": {
				Type: data.QuestionNoul,
				Instructions: data.Ask("Does `extracted_value` match the `field` as it appears in `source_text`?").
					WithField(data.Field{
						Name:        "invoice_number",
						Type:        "string",
						Description: "The identifier printed on the invoice.",
					}).
					WithExtractedValue("4471"),
			},

			// Field with a unit, judged on a scale.
			"amount_due_size": {
				Type: data.QuestionScore,
				Instructions: data.Ask("How large is the `field` value in `source_text`?").
					WithField(data.Field{
						Name:        "amount_due",
						Type:        "number",
						Unit:        "USD",
						Description: "The total the invoice asks to be paid.",
					}),
				Criteria: data.CriteriaScore{"Small", "Typical", "Unusually large"},
			},

			// State paths to weigh against each other, plus what to focus on.
			"sender_mismatch": {
				Type: data.QuestionNoul,
				Instructions: data.Ask("Does the claimed sender identity conflict with the sending domain?").
					WithCompare("ticket.sender.display_name", "ticket.sender.email").
					WithFocus("Compare the named organization with the email domain."),
			},

			// A caveat for how to judge.
			"pr_focus": {
				Type: data.QuestionScore,
				Instructions: data.Ask("How focused is this pull request description on a single change?").
					WithNote("Judge the number of independent changes, not the size of any one change."),
				Criteria: data.CriteriaScore{"One change", "A few related changes", "Unrelated changes"},
			},
		},
	}

	out, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output:
	// {
	//   "model": "jev-latest",
	//   "state": "...",
	//   "questions": {
	//     "amount_due_size": {
	//       "type": "score",
	//       "instructions": {
	//         "field": {
	//           "name": "amount_due",
	//           "type": "number",
	//           "unit": "USD",
	//           "description": "The total the invoice asks to be paid."
	//         },
	//         "question": "How large is the `field` value in `source_text`?"
	//       },
	//       "criteria": [
	//         "Small",
	//         "Typical",
	//         "Unusually large"
	//       ]
	//     },
	//     "department": {
	//       "type": "choice",
	//       "instructions": "Which team should handle this",
	//       "criteria": {
	//         "billing": "Payment or subscription issues",
	//         "technical": "Bugs or integration problems"
	//       }
	//     },
	//     "invoice_number_matches": {
	//       "type": "noul",
	//       "instructions": {
	//         "extracted_value": "4471",
	//         "field": {
	//           "name": "invoice_number",
	//           "type": "string",
	//           "description": "The identifier printed on the invoice."
	//         },
	//         "question": "Does `extracted_value` match the `field` as it appears in `source_text`?"
	//       }
	//     },
	//     "pr_focus": {
	//       "type": "score",
	//       "instructions": {
	//         "note": "Judge the number of independent changes, not the size of any one change.",
	//         "question": "How focused is this pull request description on a single change?"
	//       },
	//       "criteria": [
	//         "One change",
	//         "A few related changes",
	//         "Unrelated changes"
	//       ]
	//     },
	//     "sender_mismatch": {
	//       "type": "noul",
	//       "instructions": {
	//         "compare": [
	//           "ticket.sender.display_name",
	//           "ticket.sender.email"
	//         ],
	//         "focus": "Compare the named organization with the email domain.",
	//         "question": "Does the claimed sender identity conflict with the sending domain?"
	//       }
	//     }
	//   }
	// }
}

// ExampleInstructions_Decode shows reading structured instructions back off
// the wire, including a context key this library has no helper for.
//
//nolint:lll // the Output block must match the encoder byte for byte.
func ExampleInstructions_Decode() {
	raw := []byte(`{
		"type": "noul",
		"instructions": {
			"question": "Does ` + "`extracted_value`" + ` match the ` + "`field`" + `?",
			"field": {"name": "amount_due", "type": "number", "unit": "USD"},
			"extracted_value": 4471,
			"tolerance": {"percent": 0.5}
		}
	}`)

	var q data.Question
	if err := json.Unmarshal(raw, &q); err != nil {
		panic(err)
	}

	var f data.Field
	if err := q.Instructions.Decode("field", &f); err != nil {
		panic(err)
	}
	fmt.Printf("question: %s\n", q.Instructions.Question)
	fmt.Printf("field: %s (%s, %s)\n", f.Name, f.Type, f.Unit)

	// Keys with no typed helper are still readable and still round-trip.
	tolerance, _ := q.Instructions.Get("tolerance")
	fmt.Printf("tolerance: %s\n", tolerance)

	out, err := json.Marshal(q.Instructions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("re-encoded: %s\n", out)
}
