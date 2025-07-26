package main

import (
	"bufio"
	"fmt"
	"os"
)

func largestPrefix(heights []int64) []int64 {
	n := len(heights)
	pref := make([]int64, n+1)
	type item struct {
		pos int
		h   int64
	}
	stack := make([]item, 0)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1]
		start := i
		h := heights[i-1]
		for len(stack) > 0 && stack[len(stack)-1].h >= h {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			area := top.h * int64(i-top.pos)
			if area > pref[i] {
				pref[i] = area
			}
			start = top.pos
		}
		stack = append(stack, item{start, h})
		if h > pref[i] {
			pref[i] = h
		}
	}
	end := n + 1
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		area := top.h * int64(end-top.pos)
		if area > pref[n] {
			pref[n] = area
		}
	}
	return pref
}

func largestSuffix(heights []int64) []int64 {
	n := len(heights)
	suff := make([]int64, n+2)
	type item struct {
		pos int
		h   int64
	}
	stack := make([]item, 0)
	for i := n; i >= 1; i-- {
		suff[i] = suff[i+1]
		start := i
		h := heights[i-1]
		for len(stack) > 0 && stack[len(stack)-1].h >= h {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			area := top.h * int64(top.pos-i)
			if area > suff[i] {
				suff[i] = area
			}
			start = top.pos
		}
		stack = append(stack, item{start, h})
		if h > suff[i] {
			suff[i] = h
		}
	}
	start := 0
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		area := top.h * int64(top.pos-start)
		if area > suff[1] {
			suff[1] = area
		}
	}
	return suff
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
	}
	pref := largestPrefix(h)
	suff := largestSuffix(h)
	best := pref[n]
	for i := 1; i < n; i++ {
		sum := pref[i] + suff[i+1]
		if sum > best {
			best = sum
		}
	}
	fmt.Fprintln(out, best)
}
