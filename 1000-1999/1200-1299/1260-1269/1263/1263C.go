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
		vals := make(map[int]struct{})
		vals[0] = struct{}{}
		for i := 1; i*i <= n; i++ {
			vals[i] = struct{}{}
			vals[n/i] = struct{}{}
		}
		ans := make([]int, 0, len(vals))
		for v := range vals {
			ans = append(ans, v)
		}
		sort.Ints(ans)
		fmt.Fprintln(writer, len(ans))
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
