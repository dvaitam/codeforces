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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	coords := make([]int, n)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		coords[i-1] = a[i]
	}
	sort.Ints(coords)
	uniq := make([]int, 0, n)
	for i := 0; i < n; {
		v := coords[i]
		uniq = append(uniq, v)
		j := i + 1
		for j < n && coords[j] == v {
			j++
		}
		i = j
	}
	posMap := make(map[int]int, len(uniq))
	for idx, v := range uniq {
		posMap[v] = idx
	}
	dp := make([]int, n+1)
	last := make([]int, len(uniq))
	for i := 1; i <= n; i++ {
		dp[i] = dp[i-1]
		pos := posMap[a[i]]
		if last[pos] != 0 {
			ptr := last[pos]
			if dp[ptr-1]+1 > dp[i] {
				dp[i] = dp[ptr-1] + 1
			}
		}
		last[pos] = i
	}
	if dp[n] == 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	fmt.Fprintln(writer, dp[n])
	r := n
	i, j := n-2, n-2
	for i > -1 {
		if dp[j]+1 != dp[r] {
			j = i
			i = j - 1
			continue
		}
		for j = i; j > -1 && dp[j]+1 == dp[r]; j-- {
		}
		j++
		// print segment from j+1 to r
		fmt.Fprintln(writer, j+1, r)
		r = j
		i = j - 1
	}
}
