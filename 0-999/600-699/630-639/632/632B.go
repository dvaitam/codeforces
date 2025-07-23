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

	p := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}

	var s string
	fmt.Fscan(in, &s)

	prefixA := make([]int64, n+1)
	prefixB := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefixA[i+1] = prefixA[i]
		prefixB[i+1] = prefixB[i]
		if s[i] == 'A' {
			prefixA[i+1] += p[i]
		} else {
			prefixB[i+1] += p[i]
		}
	}

	totalA := prefixA[n]
	totalB := prefixB[n]
	ans := totalB
	for k := 0; k <= n; k++ {
		val := totalB - prefixB[k] + prefixA[k]
		if val > ans {
			ans = val
		}
		val = prefixB[n-k] + (totalA - prefixA[n-k])
		if val > ans {
			ans = val
		}
	}

	fmt.Fprintln(out, ans)
}
