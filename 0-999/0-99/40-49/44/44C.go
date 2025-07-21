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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	// count watering per day
	cnt := make([]int, n+2)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		for d := a; d <= b; d++ {
			cnt[d]++
		}
	}
	// find first mistake
	for d := 1; d <= n; d++ {
		if cnt[d] != 1 {
			fmt.Fprintln(out, d, cnt[d])
			return
		}
	}
	fmt.Fprintln(out, "OK")
}
