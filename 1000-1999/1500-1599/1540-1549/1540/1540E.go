package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	g := make([][]int, n)
	for i := 0; i < n; i++ {
		var c int
		fmt.Fscan(in, &c)
		g[i] = make([]int, c)
		for j := 0; j < c; j++ {
			fmt.Fscan(in, &g[i][j])
			g[i][j]--
		}
	}
	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var k, l, r int
			fmt.Fscan(in, &k, &l, &r)
			l--
			r--
			tmp := make([]int64, n)
			copy(tmp, a)
			for day := 0; day < k; day++ {
				stage1 := make([]int64, n)
				for i := 0; i < n; i++ {
					if tmp[i] > 0 {
						stage1[i] = tmp[i] * int64(i+1)
					} else {
						stage1[i] = tmp[i]
					}
				}
				next := make([]int64, n)
				for i := 0; i < n; i++ {
					val := stage1[i]
					for _, j := range g[i] {
						if stage1[j] > 0 {
							val += stage1[j]
						}
					}
					next[i] = val
				}
				tmp = next
			}
			sum := int64(0)
			for i := l; i <= r; i++ {
				sum = (sum + tmp[i]) % MOD
			}
			if sum < 0 {
				sum += MOD
			}
			fmt.Fprintln(out, sum)
		} else {
			var idx int
			var x int64
			fmt.Fscan(in, &idx, &x)
			a[idx-1] += x
		}
	}
}
