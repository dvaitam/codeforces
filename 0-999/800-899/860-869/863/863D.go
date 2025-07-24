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

	var n, q, m int
	if _, err := fmt.Fscan(reader, &n, &q, &m); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	t := make([]int, q)
	l := make([]int, q)
	r := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &t[i], &l[i], &r[i])
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i])
	}

	res := make([]int, m)
	for idx := 0; idx < m; idx++ {
		pos := b[idx]
		for i := q - 1; i >= 0; i-- {
			if pos < l[i] || pos > r[i] {
				continue
			}
			if t[i] == 1 {
				if pos == l[i] {
					pos = r[i]
				} else {
					pos--
				}
			} else {
				pos = l[i] + r[i] - pos
			}
		}
		res[idx] = a[pos]
	}

	for i := 0; i < m; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res[i])
	}
	fmt.Fprintln(writer)
}
