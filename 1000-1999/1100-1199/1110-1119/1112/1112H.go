package main

import (
	"bufio"
	"fmt"
	"os"
)

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var c, d int64
	fmt.Fscan(reader, &n, &c, &d)
	t := make([]int64, n)
	p := make([]byte, n)
	for i := 0; i < n; i++ {
		var pi string
		fmt.Fscan(reader, &t[i], &pi)
		if len(pi) > 0 {
			p[i] = pi[0]
		}
	}
	var finalT int64
	fmt.Fscan(reader, &finalT)
	// next event time for W and P (when they pick up letters)
	nextW, nextP := finalT, finalT
	var total int64
	// traverse events backwards
	for i := n - 1; i >= 0; i-- {
		ti := t[i]
		if p[i] == 'W' {
			// W sends to P, pick at nextP
			delta := nextP - ti
			cost := min64(d, c*delta)
			total += cost
			// update next W event
			nextW = ti
		} else {
			// P sends to W
			delta := nextW - ti
			cost := min64(d, c*delta)
			total += cost
			nextP = ti
		}
	}
	fmt.Println(total)
}
