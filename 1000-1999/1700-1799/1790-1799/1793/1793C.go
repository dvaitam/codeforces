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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		l, r := 0, n-1
		mn, mx := 1, n
		ansL, ansR := -1, -1
		for l < r {
			if a[l] != mn && a[l] != mx && a[r] != mn && a[r] != mx {
				ansL, ansR = l+1, r+1
				break
			}
			moved := false
			if a[l] == mn || a[l] == mx {
				if a[l] == mn {
					mn++
				} else {
					mx--
				}
				l++
				moved = true
			}
			if l >= r {
				break
			}
			if a[r] == mn || a[r] == mx {
				if a[r] == mn {
					mn++
				} else {
					mx--
				}
				r--
				moved = true
			}
			if !moved {
				l++
			}
		}
		if ansL == -1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ansL, ansR)
		}
	}
}
