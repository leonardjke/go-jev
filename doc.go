// Package jev is a client for the Jev judgement API.
//
// Give it a piece of state and a set of questions; it returns one structured
// answer per question, keyed the way you asked them.
//
//	client, err := jev.New(os.Getenv("JEV_API_KEY"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	resp, err := client.Request(ctx, ticket,
//		jev.Choice("department", "Which team should handle this?", map[string]string{
//			"billing":   "Payment or subscription issues",
//			"technical": "Bugs or integration problems",
//		}),
//		jev.Score("frustration", "How frustrated does the customer appear?",
//			"Calm", "Frustrated but civil", "Very angry"),
//		jev.Noul("is_urgent", "The message conveys urgency"),
//	)
//
// # Questions
//
// [Noul] asks a yes/no question, [Choice] picks one key out of a set of
// described answers, and [Score] picks one point on an ordered scale. Those
// cover the common case on their own; the With* methods on [Question] add
// instruction context to the questions that need it and leave the rest alone.
// [Question.Set] is the escape hatch for context keys with no typed helper
// yet, and unrecognized keys survive a decode and re-encode untouched.
//
// # Answers
//
// Answers come back under the keys they were asked with, in
// [github.com/leonardjke/go-jev/data.Response.Answers], alongside the
// confidence, the probability spread the answer was drawn from and, for a
// score, the legend mapping positions back to labels. Score, Noul and
// Confidence are pointers, so a returned 0 is distinguishable from "not
// returned".
//
// # Configuration
//
// [New] takes the API key plus any number of [Option] values: [WithModel],
// [WithEndpoint], [WithClient] and [WithLogger]. The key is held in a type
// that redacts itself through fmt, JSON and encoding.TextMarshaler, so
// printing a client or logging it never discloses the credential.
//
// Non-2xx responses come back as [*APIError], matchable with errors.As.
package jev
