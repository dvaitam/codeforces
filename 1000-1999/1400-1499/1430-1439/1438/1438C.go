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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, m int
	fmt.Fscan(reader, &n, &m)
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}

	for i := 0; i < n; i++ {
		z := (i + 1) % 2
		for j := 0; j < m; j++ {
			if a[i][j]%2 == z {
				z = 1 - z
				continue
			}
			a[i][j]++
			z = 1 - z
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, a[i][j])
		}
		writer.WriteByte('\n')
	}
}
