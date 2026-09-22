// Package privacy holds the types that keep credentials out of logs, error
// messages and anything else that formats or marshals a value.
package privacy

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
)

const redacted = "[REDACTED]"

// SensitiveString is a string that never prints or marshals itself. The real
// value is reachable only through an explicit string() conversion, which is
// the one place a caller has to mean it.
type SensitiveString string

var (
	_ fmt.Formatter          = (*SensitiveString)(nil)
	_ fmt.Stringer           = (*SensitiveString)(nil)
	_ json.Marshaler         = (*SensitiveString)(nil)
	_ encoding.TextMarshaler = (*SensitiveString)(nil)
)

// Format redacts every verb, not just the string ones: %d on a plain
// fmt.Stringer would print %!d(privacy.SensitiveString=<the secret>).
func (ss SensitiveString) Format(f fmt.State, _ rune) {
	_, _ = io.WriteString(f, redacted)
}

func (ss SensitiveString) MarshalJSON() ([]byte, error) {
	return json.Marshal(redacted)
}

func (ss SensitiveString) MarshalText() ([]byte, error) {
	return []byte(redacted), nil
}

func (ss SensitiveString) String() string {
	return redacted
}
