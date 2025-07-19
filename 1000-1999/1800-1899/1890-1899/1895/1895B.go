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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		total := n * 2
		a := make([]int, total)
		for i := 0; i < total; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Ints(a)
		// form pairs
		po := make([][2]int, n)
		for i := 0; i < n; i++ {
			po[i][0] = a[i]
			po[i][1] = a[total-1-i]
		}
		// compute answer
		ans := 0
		for i := 0; i+1 < n; i++ {
			ans += abs(po[i][0] - po[i+1][0])
			ans += abs(po[i][1] - po[i+1][1])
		}
		fmt.Fprintln(writer, ans)
		for i := 0; i < n; i++ {
			fmt.Fprintln(writer, po[i][0], po[i][1])
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
