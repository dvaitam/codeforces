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
	a := make([]uint64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	freq := make(map[uint64]int)
	for _, v := range a {
		freq[v]++
	}

	var good []uint64
	for v, c := range freq {
		if c > 1 {
			good = append(good, v)
		}
	}

	if len(good) == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	var ans int64
	for i := 0; i < n; i++ {
		if freq[a[i]] > 1 {
			ans += b[i]
		}
	}

	for i := 0; i < n; i++ {
		if freq[a[i]] > 1 {
			continue
		}
		for _, g := range good {
			if a[i]|g == g {
				ans += b[i]
				break
			}
		}
	}

	fmt.Fprintln(out, ans)
}
