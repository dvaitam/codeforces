package main

import (
	"bufio"
	"fmt"
	"os"
)

func minPieces(s string) int {
	n := len(s)
	segments := 0
	mixed := 0
	start := 0
	for i := 0; i < n-1; i++ {
		if s[i] == '1' && s[i+1] == '0' {
			if s[start] == '0' && s[i] == '1' {
				mixed++
			}
			segments++
			start = i + 1
		}
	}
	if s[start] == '0' && s[n-1] == '1' {
		mixed++
	}
	segments++
	if mixed > 0 {
		segments += mixed - 1
	}
	return segments
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, minPieces(s))
	}
}
