package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t string
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	// Build s' by removing all 'a' characters
	var s2 []rune
	for _, ch := range t {
		if ch != 'a' {
			s2 = append(s2, ch)
		}
	}
	// If s' is empty, original string is t
	if len(s2) == 0 {
		fmt.Println(t)
		return
	}
	// s' should appear twice in t: once from original and once as suffix
	if len(s2)%2 != 0 {
		fmt.Println(":(")
		return
	}
	half := len(s2) / 2
	prefix := string(s2[:half])
	suffix := t[len(t)-half:]
	if prefix == suffix {
		// original s is the prefix of t before the appended s'
		s := t[:len(t)-half]
		fmt.Println(s)
	} else {
		fmt.Println(":(")
	}
}
