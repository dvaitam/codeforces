package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n, k int, d int, a []int, v []int) int {
	arr := make([]int, n)
	copy(arr, a)
	limit := d
	if 2*n < limit {
		limit = 2 * n
	}
	ans := 0
	for i := 0; i < limit; i++ {
		res := 0
		for p := 0; p < n; p++ {
			if arr[p] == p+1 {
				res++
			}
		}
		if score := res + (d-i-1)/2; score > ans {
			ans = score
		}
		vi := v[i%k]
		for p := 0; p < vi; p++ {
			arr[p]++
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k, d int
		fmt.Fscan(in, &n, &k, &d)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		v := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &v[i])
		}
		fmt.Fprintln(out, solve(n, k, d, a, v))
	}
}
