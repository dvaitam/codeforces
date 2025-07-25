package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		A := make([]int64, n+1)
		B := make([]int64, n+1)
		M := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			A[i] = A[i-1] + a[i-1]
			B[i] = B[i-1] + b[i-1]
			if A[i] < B[i] {
				M[i] = A[i]
			} else {
				M[i] = B[i]
			}
		}
		r := 0
		base := int64(0)
		resets := []int{}
		success := false
		for {
			tpos := r + 1
			for tpos <= n && B[tpos]-base < k {
				tpos++
			}
			if tpos > n {
				break
			}
			if A[tpos]-base < k {
				success = true
				break
			}
			need := A[tpos] - k + 1
			j := sort.Search(len(M[:tpos]), func(i int) bool { return M[i] >= need })
			if j >= tpos {
				break
			}
			r = j
			base = M[r]
			resets = append(resets, r)
			if r >= tpos {
				tpos = r + 1
			}
		}
		if !success {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, len(resets))
			if len(resets) > 0 {
				for i, v := range resets {
					if i > 0 {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, v)
				}
				fmt.Fprintln(out)
			}
		}
	}
}
