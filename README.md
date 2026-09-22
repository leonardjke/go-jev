# go-jev

[![Go Reference](https://pkg.go.dev/badge/github.com/leonardjke/go-jev.svg)](https://pkg.go.dev/github.com/leonardjke/go-jev)
[![CI](https://github.com/leonardjke/go-jev/actions/workflows/ci.yml/badge.svg)](https://github.com/leonardjke/go-jev/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/leonardjke/go-jev)](https://goreportcard.com/report/github.com/leonardjke/go-jev)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

A Go client for the Jev judgement API.

Give it a piece of state and a set of questions; it returns one structured
answer per question, keyed the way you asked them.

## Install

```sh
go get github.com/leonardjke/go-jev
```

Requires Go 1.24 or later.

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/leonardjke/go-jev"
)

func main() {
	client, err := jev.New(os.Getenv("JEV_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	ticket := "Hi -- I was charged twice for June and nobody has replied in a week."

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
```

`state` is the text being judged. Every question carries a key of its own, and
the answers come back under those keys in `resp.Answers`.

## Question types

| Constructor | Asks for | Answer field |
| --- | --- | --- |
| `jev.Noul(key, question)` | a yes/no reading | `Answer.Noul` |
| `jev.Choice(key, question, criteria)` | one key out of `criteria` | `Answer.Choice` |
| `jev.Score(key, question, scale...)` | one point on an ordered scale | `Answer.Score` |

**Noul** is a plain check. Criteria are optional — add them with
`WithCriteria` when what counts as true needs spelling out.

```go
jev.Noul("is_urgent", "The message conveys urgency")

jev.Noul("is_spam", "The message is spam").
	WithCriteria(map[string]string{
		"true":  "Unsolicited bulk advertising",
		"false": "A genuine request, even a badly written one",
	})
```

**Choice** maps each possible answer to the description that qualifies it. The
winning key comes back in `Answer.Choice`.

```go
jev.Choice("department", "Which team should handle this?", map[string]string{
	"billing":   "Payment or subscription issues",
	"technical": "Bugs or integration problems",
})
```

**Score** takes scale labels lowest to highest. `Answer.Score` is the position
picked, and `Answer.Legend` maps positions back to their labels.

```go
jev.Score("frustration", "How frustrated does the customer appear?",
	"Calm", "Frustrated but civil", "Very angry")
```

## Reading answers

```go
a := resp.Answers["frustration"]

if a.Score != nil {
	fmt.Println(*a.Score, a.Legend) // the position, and the labels it maps to
}
if a.Confidence != nil {
	fmt.Println(*a.Confidence)
}
fmt.Println(a.Probabilities) // the spread the answer was picked from
fmt.Println(resp.Usage.InputTokens, resp.Usage.OutputTokens)
```

`Score`, `Noul` and `Confidence` are pointers because `0` is a meaningful
value for each of them and has to be told apart from "not returned".

## Question context

Questions that need more than their text take it through chained `With*`
calls. The ones that don't stay a single line.

```go
resp, err := client.Request(ctx, invoice,
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

	jev.Score("pr_focus", "How focused is this pull request on a single change?",
		"One change", "A few related changes", "Unrelated changes").
		WithNote("Judge the number of independent changes, not the size of any one change."),
)
```

| Method | Instruction key | Holds |
| --- | --- | --- |
| `WithField` | `field` | the value being judged |
| `WithExtractedValue` | `extracted_value` | the candidate to check the field against |
| `WithCompare` | `compare` | state paths to weigh against each other |
| `WithFocus` | `focus` | what the judge should look at |
| `WithNote` | `note` | a caveat for how to judge |
| `WithCriteria` | — | what makes a `Noul` answer true or false |

`Set(key, value)` is the escape hatch for any instruction key without a typed
helper yet:

```go
jev.Noul("within_tolerance", "Is `extracted_value` within tolerance of the `field`?").
	WithExtractedValue(4471).
	Set("tolerance", map[string]float64{"percent": 0.5})
```

Unrecognized keys also survive a decode and re-encode untouched, so a request
read off the wire round-trips.

## Client options

```go
client, err := jev.New(apiKey,
	jev.WithModel("jev-latest"),
	jev.WithEndpoint("https://gateway.internal/v1/systemone"),
	jev.WithClient(&http.Client{Timeout: 10 * time.Second}),
	jev.WithLogger(slog.Default()),
)
```

| Option | Default | Notes |
| --- | --- | --- |
| `WithModel` | `jev-latest` | model to judge with |
| `WithEndpoint` | `https://api.typesafe.ai/v1/systemone` | full URL the request is posted to — nothing is appended, so a proxy or compatible service needs its whole path |
| `WithClient` | `&http.Client{Timeout: 30 * time.Second}` | a nil client is ignored rather than applied |
| `WithLogger` | discards output | any `Debug(msg string, args ...any)`, so `*slog.Logger` fits as-is; a nil logger is rejected by `New` |

Responses are read through a limit of 8 MiB.

## API keys stay out of logs

The key is held as a type that redacts itself through `fmt`, `String()`, JSON
and `encoding.TextMarshaler`, for every verb rather than just the string ones.
Printing a client yields the endpoint and model and never the credential:

```go
fmt.Println(client) // Client{endpoint: "https://api.typesafe.ai/v1/systemone", model: "jev-latest"}
```

Reaching the real value takes an explicit `string()` conversion, which is the
one place a caller has to mean it.

## Errors

Non-2xx responses come back as `*jev.APIError`:

```go
resp, err := client.Request(ctx, state, questions...)

var apiErr *jev.APIError
switch {
case errors.As(err, &apiErr):
	log.Printf("status %d: %s", apiErr.StatusCode, apiErr.Body)
case err != nil:
	log.Fatal(err)
}
```

Request assembly is validated before anything is sent. An empty question key,
a duplicate key, no questions at all, or a failed `Set` value each return an
error naming the question at fault.

## The `data` package

`github.com/leonardjke/go-jev/data` holds the request and response types
directly and is usable on its own — for a server implementing this API, or for
tests that build a payload by hand. `data.Question`, `data.Instructions`,
`data.Response` and friends live there; the constructors in the root package
are a layer on top.

Criteria decode into the concrete type matching the question type, so you get
a typed value back rather than an `any` to assert on:

```go
scale, ok := data.QuestionCriteria[data.CriteriaScore](q)    // "score"
answers, ok := data.QuestionCriteria[data.CriteriaOthers](q) // "noul" / "choice"
```

## Testing

```sh
go test ./...
```

End-to-end tests are behind the `e2e` build tag and talk to the live API. They
skip unless `API_KEY` is set — copy `.env.example` to `.env` and fill it in:

```sh
API_KEY=... API_URL=... go test -tags e2e -v ./...
```

## License

Apache 2.0 — see [LICENSE](LICENSE) and [NOTICE](NOTICE).
