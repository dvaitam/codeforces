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
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimRight(line, "\n\r")

	var tokens []string
	for i := 0; i < len(line); {
		switch line[i] {
		case ' ':
			i++
		case '.', ',':
			tokens = append(tokens, line[i:i+1])
			i++
		default:
			if isLetter(line[i]) {
				j := i
				for j < len(line) && isLetter(line[j]) {
					j++
				}
				tokens = append(tokens, line[i:j])
				i = j
			} else {
				i++
			}
		}
	}

	var sb strings.Builder
	capNext := true
	for idx, tok := range tokens {
		if tok == "." || tok == "," {
			sb.WriteString(tok)
			if idx < len(tokens)-1 {
				sb.WriteByte(' ')
			}
			if tok == "." {
				capNext = true
			}
			continue
		}

		word := strings.ToLower(tok)
		if capNext && len(word) > 0 {
			sb.WriteString(strings.ToUpper(string(word[0])))
			if len(word) > 1 {
				sb.WriteString(word[1:])
			}
			capNext = false
		} else {
			sb.WriteString(word)
		}
		if idx < len(tokens)-1 && tokens[idx+1] != "." && tokens[idx+1] != "," {
			sb.WriteByte(' ')
		}
	}

	fmt.Println(sb.String())
}
