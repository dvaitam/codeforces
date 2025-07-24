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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	total := int64(0)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		total += a[i]
	}

	// Total must be even to split into equal halves
	if total%2 != 0 {
		fmt.Fprintln(out, "NO")
		return
	}
	target := total / 2

	left := make(map[int64]int)
	right := make(map[int64]int)
	for _, v := range a {
		right[v]++
	}

	prefix := int64(0)
	for i := 0; i < n; i++ {
		v := a[i]
		prefix += v
		right[v]--
		if right[v] == 0 {
			delete(right, v)
		}
		left[v]++

		if prefix == target {
			fmt.Fprintln(out, "YES")
			return
		}
		if prefix < target {
			diff := target - prefix
			if right[diff] > 0 {
				fmt.Fprintln(out, "YES")
				return
			}
		} else { // prefix > target
			diff := prefix - target
			if left[diff] > 0 {
				fmt.Fprintln(out, "YES")
				return
			}
		}
	}
	fmt.Fprintln(out, "NO")
}
