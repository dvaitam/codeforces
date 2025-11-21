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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for test := 0; test < t; test++ {
		var n, q int
		fmt.Fscan(in, &n, &q)
		var s string
		fmt.Fscan(in, &s)
		bytes := []byte(s)
		hasB := false
		for _, ch := range bytes {
			if ch == 'B' {
				hasB = true
				break
			}
		}
		for qi := 0; qi < q; qi++ {
			var a int64
			fmt.Fscan(in, &a)
			var steps int64
			if !hasB {
				steps = a
			} else {
				idx := 0
				x := a
				for x > 0 {
					if bytes[idx] == 'A' {
						x--
					} else {
						x /= 2
					}
					steps++
					idx++
					if idx == n {
						idx = 0
					}
				}
			}
			if qi > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, steps)
		}
		if test+1 < t {
			fmt.Fprintln(out)
		} else if q > 0 {
			fmt.Fprintln(out)
		}
	}
}
