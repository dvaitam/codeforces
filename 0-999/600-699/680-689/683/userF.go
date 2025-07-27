package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	// tokenize
	var tokens []string
	for i := 0; i < len(s); {
		switch s[i] {
		case ' ':
			i++
		case '.', ',':
			tokens = append(tokens, s[i:i+1])
			i++
		default:
			if isLetter(s[i]) {
				j := i
				for j < len(s) && isLetter(s[j]) {
					j++
				}
				tokens = append(tokens, s[i:j])
				i = j
			} else {
				i++
			}
		}
	}

	var out strings.Builder
	capNext := true
	prev := "" // "", "word", "punct"

	for _, tk := range tokens {
		if tk == "." || tk == "," {
			out.WriteString(tk)
			if tk == "." {
				capNext = true
			}
			prev = "punct"
			continue
		}

		if prev != "" {
			out.WriteByte(' ')
		}
		word := strings.ToLower(tk)
		if capNext && len(word) > 0 {
			word = strings.ToUpper(word[:1]) + word[1:]
		}
		out.WriteString(word)
		capNext = false
		prev = "word"
	}

	fmt.Println(out.String())
}