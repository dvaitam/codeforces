package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	var m int
	if _, err := fmt.Fscan(in, &m); err != nil {
		return
	}

	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
	}

	i, j := 0, 0
	first := true
	for i < n && j < m {
		var val int
		if a[i] <= b[j] {
			val = a[i]
			i++
		} else {
			val = b[j]
			j++
		}
		if !first {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, val)
		first = false
	}

	for i < n {
		if !first {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, a[i])
		i++
		first = false
	}

	for j < m {
		if !first {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, b[j])
		j++
		first = false
	}
}

