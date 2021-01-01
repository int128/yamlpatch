package jsonpath

import (
	"fmt"
	"strconv"
	"strings"
)

// FromJSONPointer transforms a JSON Pointer to JSON Path.
// See https://tools.ietf.org/html/rfc6901.
func FromJSONPointer(jsonPointer string) string {
	jsonPointer = strings.TrimPrefix(jsonPointer, "/")
	elements := strings.Split(jsonPointer, "/")
	var b strings.Builder
	for _, element := range elements {
		if i, err := strconv.Atoi(element); err == nil {
			_, _ = fmt.Fprintf(&b, "[%d]", i)
			continue
		}
		// unescape: https://tools.ietf.org/html/rfc6901#section-4
		element = strings.ReplaceAll(element, "~1", "/")
		element = strings.ReplaceAll(element, "~0", "~")
		// escape backslash and double-quote: https://tools.ietf.org/html/rfc6901#section-5
		element = strconv.Quote(element)
		_, _ = fmt.Fprintf(&b, `[%s]`, element)
	}
	return "$" + b.String()
}
