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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	aa := make([]int, 0, n-1)
	for _, v := range a {
		if v != 0 {
			aa = append(aa, v)
		}
	}
	bb := make([]int, 0, n-1)
	for _, v := range b {
		if v != 0 {
			bb = append(bb, v)
		}
	}
	if len(aa) != len(bb) {
		fmt.Fprintln(writer, "NO")
		return
	}
	m := len(aa)
	if m == 0 {
		fmt.Fprintln(writer, "YES")
		return
	}
	pos := make([]int, n)
	for i, v := range aa {
		pos[v] = i
	}
	shift := -1
	for i, v := range bb {
		idx := pos[v]
		d := (i - idx + m) % m
		if shift == -1 {
			shift = d
		} else if shift != d {
			fmt.Fprintln(writer, "NO")
			return
		}
	}
	fmt.Fprintln(writer, "YES")
}
