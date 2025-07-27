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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	diff := make([]int, n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
		diff[i] = a[i] - b[i]
	}

	sort.Ints(diff)

	var ans int64
	for i := 0; i < n; i++ {
		target := -diff[i]
		j := sort.Search(n, func(k int) bool { return diff[k] > target })
		if j < i+1 {
			j = i + 1
		}
		ans += int64(n - j)
	}

	fmt.Fprintln(out, ans)
}
