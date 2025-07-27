package main

import (
	"bufio"
	"fmt"
	"os"
)

// This solution implements a greedy algorithm to compute the
// minimum number of subsequence rotations required to transform
// string s into string t. The key observation is that every
// operation can fix an alternating subsequence of mismatched
// characters in one go. We process the mismatches from left to
// right and greedily append them to existing alternating chains
// or start new ones when necessary. The number of chains created
// equals the minimal number of operations.

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	// Check if the number of ones matches; otherwise impossible.
	cntS, cntT := 0, 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			cntS++
		}
		if t[i] == '1' {
			cntT++
		}
	}
	if cntS != cntT {
		fmt.Fprintln(writer, -1)
		return
	}

	// end1: number of active alternating sequences ending with '1'
	// end0: number of active alternating sequences ending with '0'
	end1, end0 := 0, 0
	operations := 0

	for i := 0; i < n; i++ {
		if s[i] == t[i] {
			continue
		}
		if s[i] == '1' { // type "10" -> represented as +1
			if end0 > 0 {
				end0--
				end1++
			} else {
				operations++
				end1++
			}
		} else { // type "01" -> represented as -1
			if end1 > 0 {
				end1--
				end0++
			} else {
				operations++
				end0++
			}
		}
	}

	fmt.Fprintln(writer, operations)
}
