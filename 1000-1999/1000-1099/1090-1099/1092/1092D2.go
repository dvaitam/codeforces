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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+2)
	maxv := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > maxv {
			maxv = a[i]
		}
	}
	n++
	a[n] = maxv
	fa := make([]int, n+1)
	for i := 1; i <= n; i++ {
		j := i - 1
		for a[i] > a[j] {
			if (i-fa[j])%2 == 0 {
				fmt.Fprintln(writer, "NO")
				return
			}
			j = fa[j]
		}
		if a[i] == a[j] {
			fa[i] = fa[j]
		} else {
			fa[i] = j
		}
	}
	fmt.Fprintln(writer, "YES")
}
