package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if n == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		first := make([]int, n+1)
		last := make([]int, n+1)
		for i := 1; i <= n; i++ {
			first[i] = -1
			last[i] = -1
		}
		for i, v := range a {
			if first[v] == -1 {
				first[v] = i
			}
			last[v] = i
		}
		vals := make([]int, 0)
		for v := 1; v <= n; v++ {
			if first[v] != -1 {
				vals = append(vals, v)
			}
		}
		if len(vals) <= 1 {
			fmt.Fprintln(writer, 0)
			continue
		}
		best := 1
		cur := 1
		for i := 1; i < len(vals); i++ {
			if last[vals[i-1]] < first[vals[i]] {
				cur++
			} else {
				if cur > best {
					best = cur
				}
				cur = 1
			}
		}
		if cur > best {
			best = cur
		}
		ans := len(vals) - best
		fmt.Fprintln(writer, ans)
	}
}
