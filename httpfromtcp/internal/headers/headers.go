package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

const (
	crlf = "\r\n"
	// single space = SP
	space           = " "
	keyValSeparator = ":"
)

var (
	crlfByte = []byte(crlf)
)

type Headers map[string]string

func New() Headers {
	return make(Headers)
}

// Get return Key value from Headers. Get is case insensitivity.
// Get will retunr false if the given key has no value.
func (h Headers) Get(key string) (string, bool) {
	val, ok := h[keyf(key)]

	return val, ok
}

func keyf(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}

// Add will insert a key value into header
// Add will ensure case insensitivity.
// If some key already has value add will concatenate the value separated by ",space".
func (h Headers) Add(key, val string) {
	key, val = keyValf(key, val)

	existingVal, ok := h[key]
	if ok {
		h[key] = fmt.Sprintf("%s, %s", existingVal, val)
		return
	}

	h[key] = val
}

func (h Headers) Delete(key string) {
	delete(h, keyf(key))
}

// Set will ovewrite a key and value if the key already exist, if not will just add new key and value.
func (h Headers) Set(key, val string) {
	key, val = keyValf(key, val)

	_, ok := h[key]
	if ok {
		delete(h, key)
	}

	h[key] = val
}

func keyValf(key, val string) (string, string) {
	keyF := keyf(key)
	valF := strings.TrimSpace(val)

	return keyF, valF
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, crlfByte)
	if idx == -1 {
		return 0, false, nil
	}

	headerText := string(data[:idx])
	if headerText == "" {
		return len(crlfByte), true, nil
	}

	sepIdx := strings.Index(headerText, keyValSeparator)
	if sepIdx == -1 {
		return 0, false, fmt.Errorf("malformed headers - not find header separator : - headers: %s", headerText)
	}

	key := headerText[:sepIdx]
	val := headerText[sepIdx+1:]

	if err := h.valiadteKey(key); err != nil {
		return 0, false, err
	}

	h.Add(key, val)

	numBytesParsed := idx + 2

	return numBytesParsed, false, nil
}

func (h Headers) valiadteKey(key string) error {
	// see https://datatracker.ietf.org/doc/html/rfc9110#name-tokens

	if strings.HasSuffix(key, space) {
		return fmt.Errorf("malformed headers key - got key ending with space - header key: %s", key)
	}

	hasInvalidChar := false
	for _, r := range key {
		if unicode.IsLetter(r) {
			continue
		}
		if unicode.IsNumber(r) {
			continue
		}

		isSpecial := specialCharsAllowed[r]
		if isSpecial {
			continue
		}

		hasInvalidChar = true
	}

	if hasInvalidChar {
		return fmt.Errorf("malformed headers - key header with not allowed char - header key: %s", key)
	}

	return nil
}

var specialCharsAllowed = map[rune]bool{
	'!':  true,
	' ':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'|':  true,
	'~':  true,
}
