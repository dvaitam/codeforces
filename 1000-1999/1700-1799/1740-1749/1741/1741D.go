package main

import (
	"bufio"
	"fmt"
	"os"
)

// solve recursively returns minimal operations to rearrange segment a so that
// it becomes sorted sequence starting from base (i.e. base, base+1, ...).
// If impossible, returns -1.
func solve(a []int, base int) int {
	n := len(a)
	if n == 1 {
		if a[0] == base {
			return 0
		}
		return -1
	}
	mid := n / 2
	left := a[:mid]
	right := a[mid:]

	// helper to check if all values in segment are within [l,r]
	check := func(seg []int, l, r int) bool {
		for _, v := range seg {
			if v < l || v > r {
				return false
			}
		}
		return true
	}

	best := -1
	// option 1: keep order
	if check(left, base, base+mid-1) && check(right, base+mid, base+n-1) {
		op1 := solve(left, base)
		if op1 != -1 {
			op2 := solve(right, base+mid)
			if op2 != -1 {
				best = op1 + op2
			}
		}
	}
	// option 2: swap children
	if check(left, base+mid, base+n-1) && check(right, base, base+mid-1) {
		op1 := solve(left, base+mid)
		if op1 != -1 {
			op2 := solve(right, base)
			if op2 != -1 {
				val := op1 + op2 + 1
				if best == -1 || val < best {
					best = val
				}
			}
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var m int
		fmt.Fscan(reader, &m)
		arr := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		res := solve(arr, 1)
		fmt.Fprintln(writer, res)
	}
}
