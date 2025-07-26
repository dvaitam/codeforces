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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		a[0] = 1 // m = 1 in this version
		for i := 1; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		sort.Ints(a)
		sort.Ints(b)

		j := 0
		match := 0
		for _, x := range a {
			for j < n && b[j] <= x {
				j++
			}
			if j == n {
				break
			}
			match++
			j++
		}

		fmt.Fprintln(writer, n-match)
	}
}
