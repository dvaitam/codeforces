package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	x int64
	y int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i], &ys[i])
	}

	cnt := make(map[pair]int64)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p := pair{xs[i] + xs[j], ys[i] + ys[j]}
			cnt[p]++
		}
	}

	var res int64
	for _, v := range cnt {
		res += v * (v - 1) / 2
	}
	fmt.Fprintln(writer, res)
}
