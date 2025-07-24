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
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Ints(a)

		idx := 0
		selected := make([][2]int, 0, n)
		for j := 1; j <= n; j++ {
			for idx < n && a[idx] < j {
				idx++
			}
			if idx == n {
				break
			}
			selected = append(selected, [2]int{a[idx], j})
			idx++
		}
		k := len(selected)
		base := 0
		for _, v := range a {
			if v > k {
				base += v - k
			}
		}
		delta := 0
		for _, p := range selected {
			ai, j := p[0], p[1]
			if ai > k {
				delta += k - j
			} else {
				delta += ai - j
			}
		}
		fmt.Fprintln(writer, base+delta)
	}
}
