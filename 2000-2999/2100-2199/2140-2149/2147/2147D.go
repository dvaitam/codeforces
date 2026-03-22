package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
		}

		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

		oddCounts := make([]int64, 0)
		for i := 0; i < n; {
			j := i + 1
			for j < n && a[j] == a[i] {
				j++
			}
			if a[i]&1 == 1 {
				oddCounts = append(oddCounts, int64(j-i))
			}
			i = j
		}

		var diff int64
		for i := len(oddCounts) - 1; i >= 0; i-- {
			if oddCounts[i] >= diff {
				diff = oddCounts[i] - diff
			} else {
				diff = diff - oddCounts[i]
			}
		}

		alice := (sum + diff) / 2
		bob := sum - alice
		fmt.Fprintln(out, alice, bob)
	}
}
