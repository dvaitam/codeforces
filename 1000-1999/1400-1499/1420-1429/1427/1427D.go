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
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var op [][]int
	// First phase: bring max and min to ends
	for l := 0; l < n/2; l++ {
		r := n - l - 1
		// find max and min indices in [l, r]
		indMax, indMin := l, l
		for i := l; i <= r; i++ {
			if a[i] > a[indMax] {
				indMax = i
			}
			if a[i] < a[indMin] {
				indMin = i
			}
		}
		mn, mx := indMin, indMax
		if mn > mx {
			mn, mx = mx, mn
		}
		curr := make([]int, 0, n)
		for i := 0; i < l; i++ {
			curr = append(curr, 1)
		}
		curr = append(curr, mn-l+1)
		if mx-mn > 1 {
			curr = append(curr, mx-mn-1)
		}
		curr = append(curr, r-mx+1)
		for i := r + 1; i < n; i++ {
			curr = append(curr, 1)
		}
		op = append(op, curr)
		perform(curr, &a)
	}
	// Second phase: fix remaining prefix
	for i := n/2 - 1; i >= 0; i-- {
		if a[i] != i+1 {
			curr := make([]int, 0, n)
			for j := 0; j <= i; j++ {
				curr = append(curr, 1)
			}
			if rem := n - 2*(i+1); rem > 0 {
				curr = append(curr, rem)
			}
			for j := 0; j <= i; j++ {
				curr = append(curr, 1)
			}
			op = append(op, curr)
			perform(curr, &a)
		}
	}
	// Output operations
	fmt.Fprintln(writer, len(op))
	for _, v := range op {
		fmt.Fprint(writer, len(v))
		for _, x := range v {
			fmt.Fprint(writer, " ", x)
		}
		fmt.Fprintln(writer)
	}
}

// perform applies a single operation: splits a into blocks of sizes curr,
// reverses the blocks, and rebuilds a.
func perform(curr []int, a *[]int) {
	var order [][]int
	j := 0
	for _, size := range curr {
		block := make([]int, size)
		copy(block, (*a)[j:j+size])
		order = append(order, block)
		j += size
	}
	// reverse block order
	for i, k := 0, len(order)-1; i < k; i, k = i+1, k-1 {
		order[i], order[k] = order[k], order[i]
	}
	// rebuild a
	var newA []int
	for _, blk := range order {
		newA = append(newA, blk...)
	}
	*a = newA
}
