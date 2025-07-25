package main

import (
	"bufio"
	"fmt"
	"os"
)

func step(a []int) bool {
	changed := false
	n := len(a)
	for i := 0; i < n-1; i++ {
		v := a[i+1] - a[i]
		if v < 0 {
			v = 0
		}
		if v != a[i+1] {
			a[i+1] = v
			changed = true
		}
	}
	v := a[0] - a[n-1]
	if v < 0 {
		v = 0
	}
	if v != a[0] {
		a[0] = v
		changed = true
	}
	return changed
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	for iter := 0; iter < 200000; iter++ {
		if !step(a) {
			break
		}
	}

	indices := make([]int, 0)
	for i, v := range a {
		if v > 0 {
			indices = append(indices, i+1)
		}
	}
	fmt.Fprintln(writer, len(indices))
	if len(indices) > 0 {
		for i, idx := range indices {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, idx)
		}
		fmt.Fprintln(writer)
	}
}

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
