package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	// compute parity of initial inversion count
	invParity := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				invParity ^= 1
			}
		}
	}

	var m int
	fmt.Fscan(in, &m)
	for ; m > 0; m-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		length := r - l + 1
		// parity of number of pairs inside the reversed segment
		pairParity := (length * (length - 1) / 2) % 2
		if pairParity == 1 {
			invParity ^= 1
		}
		if invParity == 1 {
			fmt.Fprintln(out, "odd")
		} else {
			fmt.Fprintln(out, "even")
		}
	}
}
