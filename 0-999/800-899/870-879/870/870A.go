package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It finds the smallest positive integer that contains at least one digit
// from each of the two given lists.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	a := make([]bool, 10)
	minA := 10
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		a[v] = true
		if v < minA {
			minA = v
		}
	}

	b := make([]bool, 10)
	minB := 10
	for i := 0; i < m; i++ {
		var v int
		fmt.Fscan(in, &v)
		b[v] = true
		if v < minB {
			minB = v
		}
	}

	common := 10
	for d := 1; d <= 9; d++ {
		if a[d] && b[d] {
			if d < common {
				common = d
			}
		}
	}

	if common < 10 {
		fmt.Fprintln(out, common)
		return
	}

	if minA < minB {
		fmt.Fprintf(out, "%d%d\n", minA, minB)
	} else {
		fmt.Fprintf(out, "%d%d\n", minB, minA)
	}
}
