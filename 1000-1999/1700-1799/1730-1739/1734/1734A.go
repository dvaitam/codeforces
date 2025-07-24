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
		ans := int(1e18)
		for i := 0; i+2 < n; i++ {
			diff := a[i+2] - a[i]
			if diff < ans {
				ans = diff
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
