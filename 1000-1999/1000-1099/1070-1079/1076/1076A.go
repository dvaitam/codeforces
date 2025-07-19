package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var s string
	for {
		// read n and s until EOF
		if _, err := fmt.Fscan(in, &n, &s); err != nil {
			break
		}
		// find first position where s[i] > s[i+1]
		idx := -1
		for i := 0; i+1 < len(s); i++ {
			if s[i] > s[i+1] {
				idx = i
				break
			}
		}
		// if no such position, remove last character
		if idx == -1 {
			idx = len(s) - 1
		}
		// output string without character at idx
		fmt.Println(s[:idx] + s[idx+1:])
	}
}
