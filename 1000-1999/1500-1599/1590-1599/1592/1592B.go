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
		var n, x int
		fmt.Fscan(reader, &n, &x)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		if x*2 <= n {
			fmt.Fprintln(writer, "YES")
			continue
		}

		b := append([]int(nil), a...)
		sort.Ints(b)
		left := n - x
		right := x - 1
		ok := true
		for i := left; i <= right; i++ {
			if a[i] != b[i] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
