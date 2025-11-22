package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n int, s string) string {
	cntR := 0
	for i := 0; i < n; i++ {
		if s[i] == 'R' {
			cntR++
		}
	}
	cntB := n - cntR

	// If one color appears at most once, we can mirror routes that only use that color.
	if cntR <= 1 || cntB <= 1 {
		return "YES"
	}

	// For even n, both colors appear at least twice -> impossible.
	if n%2 == 0 {
		return "NO"
	}

	hasRR, hasBB := false, false
	for i := 0; i < n; i++ {
		if s[i] == 'R' && s[(i+1)%n] == 'R' {
			hasRR = true
		}
		if s[i] == 'B' && s[(i+1)%n] == 'B' {
			hasBB = true
		}
	}

	// Both colors having a double creates vertices with disjoint incident colors.
	if hasRR && hasBB {
		return "NO"
	}

	// Odd alternating cycle works (every vertex can be a center).
	if !hasRR && !hasBB {
		return "YES"
	}

	// Exactly one color has a double. The other color appears as isolated edges.
	target := byte('B')
	if !hasRR {
		target = 'R'
	}
	pos := make([]int, 0)
	for i := 0; i < n; i++ {
		if s[i] == target {
			pos = append(pos, i)
		}
	}

	allEven := true
	for i, p := range pos {
		q := pos[(i+1)%len(pos)]
		gap := q - p - 1
		if gap < 0 {
			gap += n
		}
		if gap%2 == 1 {
			allEven = false
			break
		}
	}

	if allEven {
		return "NO"
	}
	return "YES"
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)

	out := bufio.NewWriter(os.Stdout)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n, &s)
		fmt.Fprintln(out, solve(n, s))
	}
	out.Flush()
}
