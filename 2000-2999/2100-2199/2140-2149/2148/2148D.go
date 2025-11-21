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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		var evenSum int64
		odds := make([]int64, 0)

		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			if x%2 == 0 {
				evenSum += x
			} else {
				odds = append(odds, x)
			}
		}

		if len(odds) == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		sort.Slice(odds, func(i, j int) bool {
			return odds[i] > odds[j]
		})

		var oddSum int64
		k := (len(odds) + 1) / 2
		for i := 0; i < k; i++ {
			oddSum += odds[i]
		}

		fmt.Fprintln(out, evenSum+oddSum)
	}
}
