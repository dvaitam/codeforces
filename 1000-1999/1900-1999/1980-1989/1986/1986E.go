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
		var k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		groups := make(map[int][]int)
		for _, v := range arr {
			r := v % k
			groups[r] = append(groups[r], v)
		}
		oddGroups := 0
		for _, g := range groups {
			if len(g)%2 == 1 {
				oddGroups++
			}
		}
		if (n%2 == 0 && oddGroups > 0) || (n%2 == 1 && oddGroups > 1) {
			fmt.Fprintln(writer, -1)
			continue
		}
		ops := int64(0)
		for _, g := range groups {
			sort.Ints(g)
			for i := 0; i+1 < len(g); i += 2 {
				ops += int64((g[i+1] - g[i]) / k)
			}
		}
		fmt.Fprintln(writer, ops)
	}
}
