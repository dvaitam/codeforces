package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt for contest 1697.
// We must determine whether string s can be transformed into t using swaps:
// "ab" -> "ba" (moving 'b' left over 'a') and "bc" -> "cb" (moving 'c' left over 'b').
// The relative order of 'a' and 'c' cannot change, 'a' can only move right and 'c'
// can only move left. The algorithm removes all 'b' characters and verifies the
// remaining sequence as well as positional constraints on 'a' and 'c'.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)
		if canTransform(s, t) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

func canTransform(s, t string) bool {
	n := len(s)
	sNoB := make([]byte, 0, n)
	tNoB := make([]byte, 0, n)
	posA_s := make([]int, 0)
	posA_t := make([]int, 0)
	posC_s := make([]int, 0)
	posC_t := make([]int, 0)

	for i := 0; i < n; i++ {
		if s[i] != 'b' {
			sNoB = append(sNoB, s[i])
			if s[i] == 'a' {
				posA_s = append(posA_s, i)
			} else { // 'c'
				posC_s = append(posC_s, i)
			}
		}
		if t[i] != 'b' {
			tNoB = append(tNoB, t[i])
			if t[i] == 'a' {
				posA_t = append(posA_t, i)
			} else {
				posC_t = append(posC_t, i)
			}
		}
	}

	if string(sNoB) != string(tNoB) {
		return false
	}

	for i := range posA_s {
		if posA_s[i] > posA_t[i] {
			return false
		}
	}
	for i := range posC_s {
		if posC_s[i] < posC_t[i] {
			return false
		}
	}
	return true
}
