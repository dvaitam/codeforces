package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		if canDriftForever(s) {
			fmt.Fprintln(out, -1)
			continue
		}
		fmt.Fprintln(out, maxTime(s))
	}
}

func canDriftForever(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if (s[i] == '>' || s[i] == '*') && (s[i+1] == '<' || s[i+1] == '*') {
			return true
		}
	}
	return false
}

func maxTime(s string) int {
	prefix := 0
	for prefix < len(s) && s[prefix] == '<' {
		prefix++
	}
	suffix := 0
	for suffix < len(s) && s[len(s)-1-suffix] == '>' {
		suffix++
	}
	if strings.IndexByte(s, '*') != -1 {
		return max(prefix+1, suffix+1)
	}
	return max(prefix, suffix)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
