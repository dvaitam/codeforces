package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var n, m, k int
	if _, err := fmt.Fscan(os.Stdin, &n, &m, &k); err != nil {
		return
	}
	total := n * m
	base := total / k
	extra := total % k
	coords := make([][2]int, total)
	idx := 0
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			for j := 0; j < m; j++ {
				coords[idx] = [2]int{i + 1, j + 1}
				idx++
			}
		} else {
			for j := m - 1; j >= 0; j-- {
				coords[idx] = [2]int{i + 1, j + 1}
				idx++
			}
		}
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	pos := 0
	for t := 1; t <= k; t++ {
		size := base
		if t == k {
			size += extra
		}
		fmt.Fprint(w, size)
		for i := 0; i < size; i++ {
			c := coords[pos]
			fmt.Fprint(w, " ", c[0], " ", c[1])
			pos++
		}
		fmt.Fprint(w, "\n")
	}
}
