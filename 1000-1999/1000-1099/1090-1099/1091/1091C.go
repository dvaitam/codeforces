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

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var v []int64
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			if i*i != n {
				v = append(v, i)
				v = append(v, n/i)
			} else {
				v = append(v, i)
			}
		}
	}
	var ans []int64
	for _, d := range v {
		l := n + 1 - d
		x := (l - 1) / d
		x++
		y := (1 + l) * x / 2
		ans = append(ans, y)
	}
	sort.Slice(ans, func(i, j int) bool { return ans[i] < ans[j] })
	for _, y := range ans {
		fmt.Fprintf(writer, "%d ", y)
	}
}
