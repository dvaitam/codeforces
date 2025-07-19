package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	var n int
	if _, err := fmt.Fscan(reader, &s, &n); err != nil {
		return
	}
	b := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	l := make([]int, n)
	for i := range b {
		l[i] = len(b[i])
	}
	p, mx, st := 0, -1, 0
	for i := 0; i < len(s); i++ {
		k := ok(s, b, l, i)
		if k != -1 {
			if i-p > mx {
				mx = i - p
				st = p
			}
			// advance start position
			if k+1 > p {
				p = k + 1
			}
		}
	}
	if len(s)-p > mx {
		mx = len(s) - p
		st = p
	}
	fmt.Println(mx, st)
}

// ok returns the maximum starting index of any pattern in b that ends at position k in s, or -1 if none
func ok(s string, b []string, l []int, k int) int {
	mx := -1
	for t := 0; t < len(b); t++ {
		start := k - l[t] + 1
		if start < 0 {
			continue
		}
		match := true
		for i := 0; i < l[t]; i++ {
			if s[start+i] != b[t][i] {
				match = false
				break
			}
		}
		if match && start > mx {
			mx = start
		}
	}
	return mx
}
