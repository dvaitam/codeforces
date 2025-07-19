package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	ch     [11]bool
	a      []int
	m      int
	writer *bufio.Writer
)

// dfs tries to build a sequence of length m, where x is current position and s is difference of sums
func dfs(x, s int) bool {
	if x > m {
		fmt.Fprintln(writer, "YES")
		for i := 1; i <= m; i++ {
			if i > 1 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, a[i])
		}
		fmt.Fprintln(writer)
		return true
	}
	for i := 10; i > s; i-- {
		if i != a[x-1] && ch[i] {
			a[x] = i
			if dfs(x+1, i-s) {
				return true
			}
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	fmt.Fscan(reader, &s)
	for i, c := range s {
		if c == '1' {
			ch[i+1] = true
		}
	}
	fmt.Fscan(reader, &m)
	a = make([]int, m+1)
	if !dfs(1, 0) {
		fmt.Fprintln(writer, "NO")
	}
}
