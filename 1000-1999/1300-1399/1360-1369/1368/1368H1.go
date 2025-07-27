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

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	left := make([]byte, n)
	right := make([]byte, n)
	top := make([]byte, m)
	bottom := make([]byte, m)

	fmt.Fscan(in, &left)
	fmt.Fscan(in, &right)
	fmt.Fscan(in, &top)
	fmt.Fscan(in, &bottom)

	var r, b int
	count := func(arr []byte) {
		for _, c := range arr {
			if c == 'R' {
				r++
			} else if c == 'B' {
				b++
			}
		}
	}

	count(left)
	count(right)
	count(top)
	count(bottom)

	if r < b {
		fmt.Fprintln(out, r)
	} else {
		fmt.Fprintln(out, b)
	}
}
