package main

import (
	"bufio"
	"fmt"
	"os"
)

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	// Append sentinel
	s += "_"
	for i := 1; i < len(s); i++ {
		prev := s[i-1]
		if !isVowel(prev) && prev != 'n' {
			cur := s[i]
			if !isVowel(cur) {
				fmt.Println("NO")
				return
			}
		}
	}
	fmt.Println("YES")
}
