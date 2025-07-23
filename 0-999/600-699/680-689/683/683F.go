package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n\r")

	tokens := make([]string, 0)
	n := len(input)
	for i := 0; i < n; {
		ch := input[i]
		if ch == ' ' {
			i++
			continue
		}
		if ch == '.' || ch == ',' {
			tokens = append(tokens, string(ch))
			i++
			continue
		}
		j := i
		for j < n {
			c := input[j]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				j++
			} else {
				break
			}
		}
		tokens = append(tokens, input[i:j])
		i = j
	}

	var sb strings.Builder
	sentenceStart := true
	for idx, tok := range tokens {
		if tok == "." || tok == "," {
			sb.WriteString(tok)
			if idx < len(tokens)-1 {
				sb.WriteByte(' ')
			}
			if tok == "." {
				sentenceStart = true
			}
			continue
		}
		word := strings.ToLower(tok)
		if sentenceStart && len(word) > 0 {
			sb.WriteString(strings.ToUpper(string(word[0])))
			if len(word) > 1 {
				sb.WriteString(word[1:])
			}
			sentenceStart = false
		} else {
			sb.WriteString(word)
		}
		if idx < len(tokens)-1 && tokens[idx+1] != "." && tokens[idx+1] != "," {
			sb.WriteByte(' ')
		}
	}

	fmt.Println(sb.String())
}
