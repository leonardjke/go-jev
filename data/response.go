package data

// Response is what the API answers a Request with. Answers is keyed by the
// same keys as Request.Questions.
type Response struct {
	Model   string            `json:"model"`
	Answers map[string]Answer `json:"answers"`
	Usage   Usage             `json:"usage"`
}

// Answer is one judged question. Which fields are set depends on Type:
// "choice" fills Choice, "score" fills Score and Legend, "noul" fills Noul.
// Score, Noul and Confidence are pointers because 0 is a meaningful value for
// each of them and has to be told apart from "not returned".
type Answer struct {
	Type QuestionType `json:"type"`

	// Choice is the winning key, set for "choice" questions.
	Choice string `json:"choice,omitempty"`
	// Score is the position on the scale, set for "score" questions.
	Score *float64 `json:"score,omitempty"`
	// Noul is the yes/no reading, set for "noul" questions.
	Noul *float64 `json:"noul,omitempty"`

	// Confidence is how sure the model is of this answer.
	Confidence *float64 `json:"confidence,omitempty"`
	// Legend maps a "score" answer's scale positions back to their labels.
	Legend map[string]string `json:"legend,omitempty"`
	// Probabilities is the spread the answer was picked from.
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
}

// Usage reports what the request cost.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
