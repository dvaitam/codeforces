package main

import (
	"bufio"
	"fmt"
	"os"
)

func isUniversal(s string) bool {
	n := len(s)
	for i := 0; i < n; i++ {
		j := n - 1 - i
		if s[i] == s[j] {
			continue
		}
		return s[i] < s[j]
	}
	return false
}

func hasDistinctLetters(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return true
		}
	}
	return false
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
		var n, k int
		var s string
		fmt.Fscan(in, &n, &k)
		fmt.Fscan(in, &s)

		if k == 0 {
			if isUniversal(s) {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}

		if hasDistinctLetters(s) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
