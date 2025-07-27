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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
		}
		sum2 := sum * 2
		if sum2%int64(n) != 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		target := sum2 / int64(n)
		count := make(map[int64]int64)
		var ans int64
		for _, v := range a {
			ans += count[target-v]
			count[v]++
		}
		fmt.Fprintln(out, ans)
	}
}
