package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N, M int
	var s, t string
	if _, err := fmt.Fscan(reader, &N, &M); err != nil {
		return
	}
	fmt.Fscan(reader, &s, &t)

	hasStar := false
	for i := 0; i < N; i++ {
		if s[i] == '*' {
			hasStar = true
			break
		}
	}
	if hasStar {
		// '*' can match any substring (including empty)
		if N-1 > M {
			fmt.Fprintln(writer, "NO")
			return
		}
		ok := true
		// match prefix before '*'
		for i := 0; i < N; i++ {
			if s[i] == '*' {
				break
			}
			if s[i] != t[i] {
				ok = false
			}
		}
		// match suffix after '*'
		for i := 0; i < N; i++ {
			if s[N-1-i] == '*' {
				break
			}
			if s[N-1-i] != t[M-1-i] {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	} else {
		// no wildcard, must match exactly
		if s == t {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
