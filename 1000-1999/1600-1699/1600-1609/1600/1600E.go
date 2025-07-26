package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	l    int
	r    int
	last int
}

var (
	arr  []int
	memo map[state]bool
)

func solve(l, r, last int) bool {
	s := state{l, r, last}
	if v, ok := memo[s]; ok {
		return v
	}
	lastVal := -int(1 << 60)
	if last != -1 {
		lastVal = arr[last]
	}
	win := false
	if l <= r {
		if arr[l] > lastVal {
			if !solve(l+1, r, l) {
				win = true
			}
		}
		if !win && arr[r] > lastVal {
			if !solve(l, r-1, r) {
				win = true
			}
		}
	}
	memo[s] = win
	return win
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	memo = make(map[state]bool)
	if solve(0, n-1, -1) {
		fmt.Print("Alice")
	} else {
		fmt.Print("Bob")
	}
}
