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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		a := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &a[i])
		}
		capacity := n*m - 2
		stored := make([]bool, k+1)
		count := 0
		p := 0
		ok := true
		for target := k; target >= 1 && ok; target-- {
			if stored[target] {
				stored[target] = false
				count--
			} else {
				for p < k && a[p] != target {
					stored[a[p]] = true
					count++
					if count > capacity {
						ok = false
						break
					}
					p++
				}
				if !ok {
					break
				}
				if p == k {
					ok = false
					break
				}
				p++ // deposit target directly
			}
		}
		if ok {
			fmt.Fprintln(out, "YA")
		} else {
			fmt.Fprintln(out, "TIDAK")
		}
	}
}
