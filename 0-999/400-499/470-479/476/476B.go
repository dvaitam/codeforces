package main

import (
	"fmt"
	"math"
)

var ans int64
var target int
var s, t string
var qCount int

func btr(idx int, curr int) {
	if idx == len(t) {
		if curr == target {
			ans++
		}
		return
	}
	switch t[idx] {
	case '?':
		btr(idx+1, curr+1)
		btr(idx+1, curr-1)
	case '+':
		btr(idx+1, curr+1)
	case '-':
		btr(idx+1, curr-1)
	}
}

func main() {
	// Read the original sequence and compute its net sum
	if _, err := fmt.Scan(&s); err != nil {
		return
	}
	if _, err := fmt.Scan(&t); err != nil {
		return
	}
	for _, c := range s {
		if c == '+' {
			target++
		} else if c == '-' {
			target--
		}
	}
	// Count question marks in t
	for _, c := range t {
		if c == '?' {
			qCount++
		}
	}
	// Backtrack to count matching assignments
	ans = 0
	btr(0, 0)
	// Compute probability: ans / 2^qCount
	denom := math.Pow(2, float64(qCount))
	res := float64(ans) / denom
	fmt.Printf("%.9f", res)
}
