package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	ans := n
	suffix := make(map[int]int)
	r := n
	for r > 0 {
		v := a[r-1]
		if suffix[v] == 0 {
			suffix[v]++
			r--
		} else {
			break
		}
	}
	ans = min(ans, r)

	prefix := make(map[int]int)
	for l := 0; l < n; l++ {
		v := a[l]
		if prefix[v] > 0 {
			break
		}
		prefix[v]++
		for r < n && suffix[v] > 0 {
			suffix[a[r]]--
			r++
		}
		ans = min(ans, r-l-1)
	}

	fmt.Fprintln(out, ans)
}
