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
		fmt.Fscan(in, &n)

		evens := make([]int64, 0)
		var maxOdd int64 = -1

		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			if x&1 == 0 {
				evens = append(evens, x)
			} else if x > maxOdd {
				maxOdd = x
			}
		}

		if len(evens) == 0 || maxOdd == -1 {
			fmt.Fprintln(out, 0)
			continue
		}

		sort.Slice(evens, func(i, j int) bool { return evens[i] < evens[j] })

		i, j := 0, len(evens)-1
		ops := int64(0)
		for i <= j {
			if evens[i] < maxOdd {
				maxOdd += evens[i]
				ops++
				i++
			} else {
				maxOdd += 2 * evens[j]
				ops += 2
				j--
			}
		}

		fmt.Fprintln(out, ops)
	}
}
