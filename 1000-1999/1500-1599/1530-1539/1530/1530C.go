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
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })

		prefA := make([]int, n+1)
		prefB := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefA[i+1] = prefA[i] + a[i]
			prefB[i+1] = prefB[i] + b[i]
		}

		extra := 0
		for {
			k := n + extra
			take := k - k/4
			var sumA int
			if take <= extra {
				sumA = take * 100
			} else {
				rem := take - extra
				if rem > n {
					rem = n
				}
				sumA = extra*100 + prefA[rem]
			}
			var sumB int
			if take > n {
				sumB = prefB[n]
			} else {
				sumB = prefB[take]
			}
			if sumA >= sumB {
				fmt.Fprintln(out, extra)
				break
			}
			extra++
		}
	}
}
