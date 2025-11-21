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

	var n, k, q int
	if _, err := fmt.Fscan(in, &n, &k, &q); err != nil {
		return
	}

	cols := make([][]int, k)
	for j := 0; j < k; j++ {
		cols[j] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			var x int
			fmt.Fscan(in, &x)
			if i == 0 {
				cols[j][i] = x
			} else {
				cols[j][i] = cols[j][i-1] | x
			}
		}
	}

	for ; q > 0; q-- {
		var m int
		fmt.Fscan(in, &m)
		L, R := 1, n
		ok := true
		for ; m > 0; m-- {
			var r int
			var op string
			var c int
			fmt.Fscan(in, &r, &op, &c)
			if !ok {
				continue
			}
			r--
			arr := cols[r]
			if op == ">" {
				idx := sort.Search(n, func(i int) bool { return arr[i] > c })
				if idx == n {
					ok = false
				} else {
					if idx+1 > L {
						L = idx + 1
					}
				}
			} else { // op == "<"
				idx := sort.Search(n, func(i int) bool { return arr[i] >= c })
				if idx == 0 {
					ok = false
				} else {
					if idx < R {
						R = idx
					}
				}
			}
			if L > R {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(out, L)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
