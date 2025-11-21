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
	fmt.Fscan(in, &n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	mp := make(map[int64]int64)
	var ans int64
	for i := 0; i < n; i++ {
		key := b[i] - int64(i)
		es := mp[key] + b[i]
		if es > ans {
			ans = es
		}
		mp[key] = es
	}
	fmt.Fprintln(out, ans)
}
