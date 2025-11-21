package main

import (
	"bufio"
	"fmt"
	"os"
)

func minCost(s string) int {
	n := len(s)

	// Base cost: presses plus number of moves (transitions including initial)
	moves := 0
	prev := byte('0')
	for i := 0; i < n; i++ {
		if s[i] != prev {
			moves++
		}
		prev = s[i]
	}
	base := n + moves

	count01, count10 := 0, 0
	for i := 0; i+1 < n; i++ {
		if s[i] == '0' && s[i+1] == '1' {
			count01++
		}
		if s[i] == '1' && s[i+1] == '0' {
			count10++
		}
	}

	delta := 0
	if count01 >= 2 || count10 >= 2 || (s[0] == '1' && count01 >= 1) {
		delta = -2
	} else {
		has0 := make([]bool, n+1)
		has1 := make([]bool, n+1)
		for i := n - 1; i >= 0; i-- {
			has0[i] = has0[i+1]
			has1[i] = has1[i+1]
			if s[i] == '0' {
				has0[i] = true
			} else {
				has1[i] = true
			}
		}

		pref0 := make([]bool, n)
		pref1 := make([]bool, n)
		seen0, seen1 := false, false
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				seen0 = true
			} else {
				seen1 = true
			}
			pref0[i] = seen0
			pref1[i] = seen1
		}

		reduce := false
		if n > 1 && s[0] == '1' && has0[1] {
			reduce = true
		}
		if !reduce {
			for i := 1; i < n; i++ {
				if s[i-1] != s[i] {
					if s[i-1] == '0' && has0[i] {
						reduce = true
						break
					}
					if s[i-1] == '1' && has1[i] {
						reduce = true
						break
					}
				}
			}
		}
		if !reduce {
			for i := 0; i+1 < n; i++ {
				if s[i] != s[i+1] {
					if s[i+1] == '0' && pref0[i] {
						reduce = true
						break
					}
					if s[i+1] == '1' && pref1[i] {
						reduce = true
						break
					}
				}
			}
		}
		if reduce {
			delta = -1
		}
	}

	return base + delta
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, minCost(s))
	}
}

