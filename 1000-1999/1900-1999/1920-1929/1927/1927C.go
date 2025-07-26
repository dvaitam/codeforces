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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		presentA := make([]bool, k+1)
		presentB := make([]bool, k+1)
		cntA, cntB := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x <= k && !presentA[x] {
				presentA[x] = true
				cntA++
			}
		}
		for i := 0; i < m; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x <= k && !presentB[x] {
				presentB[x] = true
				cntB++
			}
		}
		ok := cntA >= k/2 && cntB >= k/2
		if ok {
			for i := 1; i <= k; i++ {
				if !presentA[i] && !presentB[i] {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
