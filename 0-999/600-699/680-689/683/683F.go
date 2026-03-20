package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	s := strings.TrimRight(string(data), "\r\n")

	var b strings.Builder
	firstWord := true
	capNext := true

	for i := 0; i < len(s); {
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			break
		}

		if isLetter(s[i]) {
			start := i
			for i < len(s) && isLetter(s[i]) {
				i++
			}
			word := []byte(strings.ToLower(s[start:i]))
			if capNext && len(word) > 0 && 'a' <= word[0] && word[0] <= 'z' {
				word[0] = word[0] - 'a' + 'A'
			}
			if !firstWord {
				b.WriteByte(' ')
			}
			b.Write(word)
			firstWord = false
			capNext = false
		} else if s[i] == '.' || s[i] == ',' {
			if !firstWord {
				b.WriteByte(s[i])
			}
			if s[i] == '.' {
				capNext = true
			}
			i++
		} else {
			i++
		}
	}

	fmt.Print(b.String())
}
