package main

import (
	"bufio"
	"fmt"
	"os"
)

func add(bit []int64, idx int, val int64) {
	for idx < len(bit) {
		bit[idx] += val
		idx += idx & -idx
	}
}

func sum(bit []int64, idx int) int64 {
	var res int64
	for idx > 0 {
		res += bit[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		bit := make([]int64, n+2)
		var ans int64
		for _, v := range arr {
			ans += sum(bit, n) - sum(bit, v-1)
			add(bit, v, 1)
		}
		fmt.Fprintln(out, ans)
	}
}
