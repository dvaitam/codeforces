package main

import (
	"bufio"
	"fmt"
	"os"
)

func calc(h []int64, H int64) int64 {
	var ones, twos int64
	for _, x := range h {
		diff := H - x
		if diff > 0 {
			ones += diff % 2
			twos += diff / 2
		}
	}
	ans := twos * 2
	if tmp := ones*2 - 1; tmp > ans {
		ans = tmp
	}
	return ans
}

func solve(h []int64) int64 {
	mx := h[0]
	for _, v := range h {
		if v > mx {
			mx = v
		}
	}
	d1 := calc(h, mx)
	d2 := calc(h, mx+1)
	if d1 < d2 {
		return d1
	}
	return d2
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
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}
		fmt.Fprintln(out, solve(h))
	}
}
