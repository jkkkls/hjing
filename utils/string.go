package utils

import (
	"bytes"
)

// CloseString ...
func CloseString(s string) string {
	if len(s) == 0 {
		return ""
	}

	var (
		buff bytes.Buffer
		u    int
	)

	for i := 0; i < len(s); i++ {
		b := s[i]
		if (b < '0' || b > '9') && (b < 'A' || b > 'Z') && (b < 'a' || b > 'z') && b != '_' {
			return ""
		}

		if u == 1 || i == 0 {
			if b >= 'a' && b <= 'z' {
				b -= 'a' - 'A'
			}
			u = 0
		} else if b == '_' {
			u = 1
			continue
		}

		buff.WriteByte(b)
	}

	return buff.String()
}

// ExpandString ...
func ExpandString(s string) string {
	if len(s) == 0 {
		return ""
	}

	var buff bytes.Buffer

	for i := 0; i < len(s); i++ {
		b := s[i]
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
			if i != 0 {
				buff.WriteByte('_')
			}
			buff.WriteByte(b)
		} else {
			buff.WriteByte(b)
		}
	}

	return buff.String()
}
