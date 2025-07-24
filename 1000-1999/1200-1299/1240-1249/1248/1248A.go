package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// For each test case we are given two sets of lines:
//
//	y = x + p_i and y = -x + q_j. Two such lines
//	intersect at ((q_j - p_i)/2, (q_j + p_i)/2).
//	The intersection point has integer coordinates iff
//	p_i and q_j have the same parity. Therefore we
//	simply count how many p_i are even/odd and the
//	same for q_j, and sum the products of counts with
//	the same parity.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		evenP, oddP := 0, 0
		for i := 0; i < n; i++ {
			var p int
			fmt.Fscan(in, &p)
			if p%2 == 0 {
				evenP++
			} else {
				oddP++
			}
		}

		var m int
		fmt.Fscan(in, &m)
		evenQ, oddQ := 0, 0
		for i := 0; i < m; i++ {
			var q int
			fmt.Fscan(in, &q)
			if q%2 == 0 {
				evenQ++
			} else {
				oddQ++
			}
		}

		ans := int64(evenP*evenQ + oddP*oddQ)
		fmt.Fprintln(out, ans)
	}
}
