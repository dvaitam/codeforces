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
	// create 1-based matrix a of size (n+1) x (n+1)
	a := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		a[i] = make([]int, n+1)
	}
	// fill rows 2 to n-1
	for i := 2; i <= n-1; i++ {
		now := a[i-1][1]
		for j := 1; j <= i-1; j++ {
			now = now%(n-1) + 1
			a[i][j] = now
		}
	}
	// fill last row if applicable
	if n >= 2 {
		a[n][1] = n - 1
		for j := 2; j <= n-1; j++ {
			a[n][j] = (a[n][j-1]%(n-1)+1)%(n-1) + 1
		}
	}
	// output matrix
	for i := 1; i <= n; i++ {
		// first elements a[i][1..i]
		for j := 1; j <= i; j++ {
			fmt.Fprint(writer, a[i][j], " ")
		}
		// remaining from column i
		for j := i + 1; j <= n; j++ {
			fmt.Fprint(writer, a[j][i], " ")
		}
		fmt.Fprintln(writer)
	}
}
