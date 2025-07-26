package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		exists := make([]bool, k+1)
		pmax := make([]int, n)
		cur := 0
		for i := 0; i < n; i++ {
			v := a[i]
			if v <= k {
				exists[v] = true
			}
			if v > cur {
				cur = v
			}
			pmax[i] = cur
		}

		smax := make([]int, n)
		cur = 0
		for i := n - 1; i >= 0; i-- {
			if a[i] > cur {
				cur = a[i]
			}
			smax[i] = cur
		}

		L := make([]int, k+1)
		idx := 0
		for c := 1; c <= k; c++ {
			for idx < n && pmax[idx] < c {
				idx++
			}
			if idx < n {
				L[c] = idx
			} else {
				L[c] = n
			}
		}

		R := make([]int, k+1)
		idx = n - 1
		for c := 1; c <= k; c++ {
			for idx >= 0 && smax[idx] < c {
				idx--
			}
			if idx >= 0 {
				R[c] = idx
			} else {
				R[c] = -1
			}
		}

		for c := 1; c <= k; c++ {
			if !exists[c] {
				if c == k {
					fmt.Fprintln(writer, 0)
				} else {
					fmt.Fprint(writer, 0, " ")
				}
				continue
			}
			res := 2 * (R[c] - L[c] + 1)
			if c == k {
				fmt.Fprintln(writer, res)
			} else {
				fmt.Fprint(writer, res, " ")
			}
		}
	}
}
