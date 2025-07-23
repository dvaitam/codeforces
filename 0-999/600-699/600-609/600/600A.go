package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	if len(s) > 1 && s[0] == '0' {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func format(words []string) string {
	if len(words) == 0 {
		return "-"
	}
	return "\"" + strings.Join(words, ",") + "\""
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)

	var current strings.Builder
	var words []string
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ',' || c == ';' {
			words = append(words, current.String())
			current.Reset()
		} else {
			current.WriteByte(c)
		}
	}
	words = append(words, current.String())

	var nums, others []string
	for _, w := range words {
		if isNumber(w) {
			nums = append(nums, w)
		} else {
			others = append(others, w)
		}
	}

	fmt.Println(format(nums))
	fmt.Println(format(others))
}
