package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	var x int64
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + a[i-1]
	}

	var ans int64
	var lastDays int64
	for k := 1; k <= n; k++ {
		if prefix[k] > x {
			break
		}
		days := (x-prefix[k])/int64(k) + 1
		ans += days - lastDays
		lastDays = days
	}
	fmt.Fprintln(writer, ans)
}
