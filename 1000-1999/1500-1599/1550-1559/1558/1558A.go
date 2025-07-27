package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func possible(a, b int) []int {
	n := a + b
	set := make(map[int]struct{})
	for start := 0; start < 2; start++ {
		var serveA, serveB int
		if start == 0 {
			serveA = (n + 1) / 2
			serveB = n / 2
		} else {
			serveA = n / 2
			serveB = (n + 1) / 2
		}
		diff := a - serveA
		L := max(0, -diff)
		U := min(serveA, serveB-diff)
		if L > U {
			continue
		}
		for w := L; w <= U; w++ {
			k := 2*w + diff
			set[k] = struct{}{}
		}
	}
	res := make([]int, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	sort.Ints(res)
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		ans := possible(a, b)
		fmt.Fprintln(writer, len(ans))
		for i, v := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
