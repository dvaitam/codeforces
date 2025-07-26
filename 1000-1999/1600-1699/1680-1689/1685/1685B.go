package main

import (
	"bufio"
	"fmt"
	"os"
)

// solve tries to decide if string s can be split into a words "A", b words "B",
// c words "AB" and d words "BA". The approach implemented here only checks some
// necessary conditions using greedy pair counts. It may not cover all edge
// cases of the original problem.
func solve(a, b, c, d int, s string) bool {
	n := len(s)
	countA := 0
	for i := 0; i < n; i++ {
		if s[i] == 'A' {
			countA++
		}
	}
	countB := n - countA
	if countA != a+c+d || countB != b+c+d {
		return false
	}

	// maximum number of disjoint pairs of form "AB" greedy
	maxAB := 0
	for i := 0; i+1 < n; {
		if s[i] == 'A' && s[i+1] == 'B' {
			maxAB++
			i += 2
		} else {
			i++
		}
	}

	// maximum number of disjoint pairs of form "BA" greedy
	maxBA := 0
	for i := 0; i+1 < n; {
		if s[i] == 'B' && s[i+1] == 'A' {
			maxBA++
			i += 2
		} else {
			i++
		}
	}

	// maximum number of disjoint pairs of either form
	maxPairs := 0
	for i := 0; i+1 < n; {
		if s[i] != s[i+1] {
			maxPairs++
			i += 2
		} else {
			i++
		}
	}

	if c+d > maxPairs || c > maxAB || d > maxBA {
		return false
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	t := 0
	fmt.Fscan(in, &t)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var a, b, c, d int
		var s string
		fmt.Fscan(in, &a, &b, &c, &d)
		fmt.Fscan(in, &s)
		if solve(a, b, c, d, s) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
