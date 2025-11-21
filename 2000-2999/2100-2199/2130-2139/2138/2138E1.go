package main

import (
	"bufio"
	"fmt"
	"os"
)

func fibSequence(limit int) []int {
	fib := []int{1, 1}
	for fib[len(fib)-1] < limit {
		next := fib[len(fib)-1] + fib[len(fib)-2]
		fib = append(fib, next)
	}
	return fib
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fib := fibSequence(1e7 + 5)
	mats := make([][][]int, len(fib))

	for idx := range fib {
		size := idx + 1
		t := make([][]int, size)
		for i := range t {
			t[i] = make([]int, size)
		}
		if size == 0 {
			continue
		}
		for i := 0; i < size; i++ {
			t[i][i] = 1
		}
		for i := 0; i+1 < size; i++ {
			if i%2 == 0 {
				t[i][i+1] = 1
				t[i+1][i] = 1
			} else {
				t[i][i+1] = -1
				t[i+1][i] = 1
			}
		}
		mats[idx] = t
	}

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x int
		fmt.Fscan(in, &x)
		if x == 0 {
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, 0)
			continue
		}
		idx := 0
		for idx < len(fib) && fib[idx] != x {
			idx++
		}
		if idx == len(fib) {
			idx = len(fib) - 1
		}
		m := mats[idx]
		size := len(m)
		fmt.Fprintln(out, size)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				fmt.Fprint(out, m[i][j])
				if j+1 < size {
					fmt.Fprint(out, " ")
				}
			}
			fmt.Fprintln(out)
		}
	}
}
