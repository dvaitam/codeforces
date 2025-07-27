package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAX = 100000

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	A := make([]int, n)
	B := make([]int, n)
	freqA := make([]int, MAX+2)
	freqB := make([]int, MAX+2)

	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
		freqA[A[i]]++
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &B[i])
		freqB[B[i]]++
	}

	var Q int
	fmt.Fscan(in, &Q)

	prefixA := make([]int, MAX+2)
	prefixB := make([]int, MAX+2)
	scores := make([]int, n)

	for ; Q > 0; Q-- {
		var t int
		fmt.Fscan(in, &t)
		switch t {
		case 1:
			var k int
			var c string
			fmt.Fscan(in, &k, &c)
			k--
			freqA[A[k]]--
			if c == "+" {
				A[k]++
			} else {
				A[k]--
			}
			freqA[A[k]]++
		case 2:
			var k int
			var c string
			fmt.Fscan(in, &k, &c)
			k--
			freqB[B[k]]--
			if c == "+" {
				B[k]++
			} else {
				B[k]--
			}
			freqB[B[k]]++
		case 3:
			var k int
			fmt.Fscan(in, &k)
			k--
			// build prefixGreater arrays
			g := 0
			for v := MAX; v >= 0; v-- {
				prefixA[v] = g
				g += freqA[v]
			}
			g = 0
			for v := MAX; v >= 0; v-- {
				prefixB[v] = g
				g += freqB[v]
			}
			for i := 0; i < n; i++ {
				ra := 1 + prefixA[A[i]]
				rb := 1 + prefixB[B[i]]
				scores[i] = ra + rb
			}
			target := scores[k]
			count := 0
			for i := 0; i < n; i++ {
				if scores[i] < target {
					count++
				}
			}
			fmt.Fprintln(out, count+1)
		}
	}
}
