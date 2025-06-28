package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

const allowedNameChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	value, exists := h[key]
	if !exists {
		return ""
	}
	return value
}

func (h Headers) put(key string, value string) {
	key = strings.ToLower(key)
	oldValue, exists := h[key]
	if exists {
		h[key] = oldValue + ", " + value
		return
	}
	h[key] = value
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	done = false
	n = 0
	err = nil
	dataString := string(data)
	if !strings.Contains(dataString, CRLF) {
		return
	}
	if dataString[:2] == CRLF {
		n = 2
		done = true
		return
	}
	splitIndex := 0
	for i, r := range dataString {
		if r == ' ' {
			err = fmt.Errorf("Unexpeted whitespace in field name")
			return
		}
		if r == ':' {
			splitIndex = i
			break
		}
		if !strings.Contains(allowedNameChars, string(r)) {
			err = fmt.Errorf("Invalid character in field name")
			return
		}
	}
	fmt.Println(dataString)
	fmt.Printf("splitIndex = %d\n", splitIndex)
	if splitIndex == 0 {
		// fmt.Printf("Bad header: %s", dataString)
		if strings.Contains(dataString, ":") {
			err = fmt.Errorf("Missing field name")
			return
		}
		err = fmt.Errorf("Expected ':' in header not found")
		return
	}
	if dataString[splitIndex+1:splitIndex+3] == CRLF {
		err = fmt.Errorf("Missing field value in header")
		return
	}
	fieldName := strings.ToLower(dataString[:splitIndex])
	fieldValue := dataString[splitIndex+1:]
	valueParts := strings.Split(fieldValue, CRLF)
	fieldValue = valueParts[0]
	n = 3 + len(fieldValue) + len(fieldName)
	fieldValue = strings.TrimSpace(fieldValue)
	h.put(fieldName, fieldValue)
	return
}
