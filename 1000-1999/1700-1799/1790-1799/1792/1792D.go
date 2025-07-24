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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		perms := make([][]int, n)
		for i := 0; i < n; i++ {
			arr := make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &arr[j])
			}
			perms[i] = arr
		}

		prefix := make(map[int64]struct{}, n*m+1)
		for _, p := range perms {
			inv := make([]int, m)
			for idx, val := range p {
				inv[val-1] = idx + 1
			}
			key := int64(0)
			for k := 0; k < m; k++ {
				key = key*11 + int64(inv[k])
				prefix[key] = struct{}{}
			}
		}

		for _, p := range perms {
			key := int64(0)
			best := 0
			for k := 0; k < m; k++ {
				key = key*11 + int64(p[k])
				if _, ok := prefix[key]; ok {
					best = k + 1
				}
			}
			fmt.Fprint(out, best, " ")
		}
		fmt.Fprintln(out)
	}
}
