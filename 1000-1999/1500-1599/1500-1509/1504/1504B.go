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
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)
		balanced := make([]bool, n)
		diff := 0
		for i := 0; i < n; i++ {
			if a[i] == '1' {
				diff++
			} else {
				diff--
			}
			if diff == 0 {
				balanced[i] = true
			}
		}
		flip := false
		possible := true
		for i := n - 1; i >= 0; i-- {
			ch := a[i]
			if flip {
				if ch == '1' {
					ch = '0'
				} else {
					ch = '1'
				}
			}
			if ch == b[i] {
				continue
			}
			if !balanced[i] {
				possible = false
				break
			}
			flip = !flip
		}
		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
