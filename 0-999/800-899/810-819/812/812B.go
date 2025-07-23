package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	type floor struct{ L, R int }
	floors := make([]floor, n)
	// read floors top to bottom, store bottom-first
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		l := m + 1
		r := 0
		for j := 0; j < m+2; j++ {
			if s[j] == '1' {
				if l > j {
					l = j
				}
				r = j
			}
		}
		floors[n-1-i] = floor{l, r}
	}

	highest := -1
	for i := n - 1; i >= 0; i-- {
		if floors[i].R != 0 {
			highest = i
			break
		}
	}
	if highest == -1 {
		fmt.Println(0)
		return
	}

	width := m + 1
	const INF = int(1e9)
	dpLeft, dpRight := 0, INF

	for i := 0; i < highest; i++ {
		l, r := floors[i].L, floors[i].R
		var newLeft, newRight int
		if r == 0 {
			newLeft = min(dpLeft, dpRight+width) + 1
			newRight = min(dpRight, dpLeft+width) + 1
		} else {
			newLeft = min(dpLeft+2*r, dpRight+width) + 1
			newRight = min(dpRight+2*(width-l), dpLeft+width) + 1
		}
		dpLeft, dpRight = newLeft, newRight
	}

	// final floor
	l, r := floors[highest].L, floors[highest].R
	ans := min(dpLeft+r, dpRight+width-l)
	fmt.Println(ans)
}
