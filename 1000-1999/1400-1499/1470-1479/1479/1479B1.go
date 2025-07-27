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

	last := [2]int{-1, -1}
	seg := [2]int{0, 0}
	for i := 0; i < n; i++ {
		x := a[i]
		next := -2
		if i+1 < n {
			next = a[i+1]
		}
		if last[0] == x && last[1] == x {
			last[0] = x
		} else if last[0] == x {
			if last[1] != x {
				seg[1]++
			}
			last[1] = x
		} else if last[1] == x {
			if last[0] != x {
				seg[0]++
			}
			last[0] = x
		} else {
			if next == last[0] {
				seg[0]++
				last[0] = x
			} else if next == last[1] {
				seg[1]++
				last[1] = x
			} else {
				if seg[0] <= seg[1] {
					seg[0]++
					last[0] = x
				} else {
					seg[1]++
					last[1] = x
				}
			}
		}
	}
	fmt.Fprintln(out, seg[0]+seg[1])
}
