package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		username := ""
		for i := 1; i < len(s); i++ {
			prefix := s[:i]
			suffix := s[i:]
			if suffix[0] == '0' {
				continue
			}
			validSuffix := true
			for _, c := range suffix {
				if !unicode.IsDigit(c) {
					validSuffix = false
					break
				}
			}
			if !validSuffix {
				continue
			}
			hasLetter := false
			for _, c := range prefix {
				if unicode.IsLetter(c) {
					hasLetter = true
					break
				}
			}
			if !hasLetter {
				continue
			}
			username = prefix
			break
		}
		if username == "" {
			username = s
		}
		fmt.Fprintln(out, username)
	}
}
