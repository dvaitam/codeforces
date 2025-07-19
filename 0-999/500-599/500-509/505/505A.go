package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	// try inserting a letter at every position
	for ch := 'a'; ch <= 'z'; ch++ {
		for pos := 0; pos <= len(s); pos++ {
			// build candidate
			t := s[:pos] + string(ch) + s[pos:]
			if isPalindrome(t) {
				fmt.Print(t)
				return
			}
		}
	}
	fmt.Print("NA")
}
