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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := range arr {
		fmt.Fscan(in, &arr[i])
	}

	const maxX = 512
	inf := 1000
	dp := make([]int, maxX)
	tmp := make([]int, maxX)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = -1

	for _, v := range arr {
		copy(tmp, dp)
		for x := 0; x < maxX; x++ {
			if dp[x] < v {
				nx := x ^ v
				if v < tmp[nx] {
					tmp[nx] = v
				}
			}
		}
		dp, tmp = tmp, dp
	}

	res := make([]int, 0)
	for x := 0; x < maxX; x++ {
		if dp[x] < inf {
			res = append(res, x)
		}
	}

	fmt.Fprintln(out, len(res))
	for i, val := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, val)
	}
	fmt.Fprintln(out)
}
