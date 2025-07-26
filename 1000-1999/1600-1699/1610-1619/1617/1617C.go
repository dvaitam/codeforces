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
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		used := make([]bool, n+1)
		extras := make([]int, 0)
		for _, v := range a {
			if v >= 1 && v <= n && !used[v] {
				used[v] = true
			} else {
				extras = append(extras, v)
			}
		}
		sort.Ints(extras)
		missing := make([]int, 0)
		for i := 1; i <= n; i++ {
			if !used[i] {
				missing = append(missing, i)
			}
		}
		if len(extras) != len(missing) {
			fmt.Fprintln(writer, -1)
			continue
		}
		possible := true
		for i, m := range missing {
			if extras[i] <= 2*m {
				possible = false
				break
			}
		}
		if possible {
			fmt.Fprintln(writer, len(extras))
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
