package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// Given a string s of 'L' and 'R', the frog starts at position 0
// and wants to reach position len(s)+1. It can only move right from
// position 0 and from any position i with s[i-1]=='R'. Positions with
// 'L' force jumps to the left. To minimize the maximal jump distance d,
// we avoid landing on 'L'. The minimal d equals the largest gap between
// consecutive 'R' cells when we consider positions 0 and n+1 as 'R'.
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
		pos := []int{0}
		for i, ch := range s {
			if ch == 'R' {
				pos = append(pos, i+1)
			}
		}
		pos = append(pos, len(s)+1)
		maxGap := 0
		for i := 1; i < len(pos); i++ {
			gap := pos[i] - pos[i-1]
			if gap > maxGap {
				maxGap = gap
			}
		}
		fmt.Fprintln(out, maxGap)
	}
}
