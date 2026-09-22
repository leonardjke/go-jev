package jev_test

import (
	"context"
	"fmt"
	"log"

	"github.com/leonardjke/go-jev"
)

// Example_simple is the everyday shape: build the client once, then ask.
// Questions that need no context are a single call each.
func Example_simple() {
	client, err := jev.New("...")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Request(context.Background(), ticket,
		jev.Choice("department", "Which team should handle this?", map[string]string{
			"billing":   "Payment or subscription issues",
			"technical": "Bugs or integration problems",
		}),
		jev.Score("frustration", "How frustrated does the customer appear?",
			"Calm", "Frustrated but civil", "Very angry"),
		jev.Noul("is_urgent", "The message conveys urgency"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Answers["department"].Choice)
}

// Example_context adds instruction context to the questions that need it. The
// call is the same one as above — only the questions carrying extra context
// grow, and the plain ones stay plain.
func Example_context() {
	client, err := jev.New("...")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Request(context.Background(), invoice,
		jev.Noul("invoice_number_matches",
			"Does `extracted_value` match the `field` as it appears in `source_text`?").
			WithField(jev.Field{
				Name:        "invoice_number",
				Type:        "string",
				Description: "The identifier printed on the invoice.",
			}).
			WithExtractedValue("4471"),

		jev.Score("amount_due_size", "How large is the `field` value in `source_text`?",
			"Small", "Typical", "Unusually large").
			WithField(jev.Field{Name: "amount_due", Type: "number", Unit: "USD"}),

		jev.Noul("sender_mismatch", "Does the claimed sender identity conflict with the sending domain?").
			WithCompare("ticket.sender.display_name", "ticket.sender.email").
			WithFocus("Compare the named organization with the email domain."),

		// Set is the escape hatch for instruction keys with no typed helper.
		jev.Noul("within_tolerance", "Is `extracted_value` within tolerance of the `field`?").
			WithExtractedValue(4471).
			Set("tolerance", map[string]float64{"percent": 0.5}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Answers["invoice_number_matches"].Noul)
}

var (
	ticket  = "..."
	invoice = "..."
)
