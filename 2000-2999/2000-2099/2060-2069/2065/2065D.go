package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type arrInfo struct {
	sum   int64
	score int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		data := make([]arrInfo, n)
		for i := 0; i < n; i++ {
			var sum int64
			var score int64
			for j := 0; j < m; j++ {
				var val int64
				fmt.Fscan(in, &val)
				sum += val
				coef := int64(m - j)
				score += coef * val
			}
			data[i] = arrInfo{sum: sum, score: score}
		}
		sort.Slice(data, func(i, j int) bool {
			if data[i].sum == data[j].sum {
				return data[i].score > data[j].score
			}
			return data[i].sum > data[j].sum
		})

		var ans int64
		for _, arr := range data {
			ans += arr.score
		}
		for i, arr := range data {
			ans += int64(m) * arr.sum * int64(n-1-i)
		}

		fmt.Fprintln(out, ans)
	}
}
