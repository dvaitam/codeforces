package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k, x int64
	if _, err := fmt.Fscan(in, &n, &k, &x); err != nil {
		return
	}
	levels := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &levels[i])
	}

	sort.Slice(levels, func(i, j int) bool { return levels[i] < levels[j] })

	gaps := make([]int64, 0)
	for i := 1; i < n; i++ {
		diff := levels[i] - levels[i-1]
		if diff > x {
			gaps = append(gaps, (diff-1)/x)
		}
	}

	sort.Slice(gaps, func(i, j int) bool { return gaps[i] < gaps[j] })

	groups := len(gaps) + 1
	for _, cost := range gaps {
		if k >= cost {
			k -= cost
			groups--
		} else {
			break
		}
	}

	fmt.Fprintln(out, groups)
}
