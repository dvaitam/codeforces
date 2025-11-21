package main

import (
	"bufio"
	"fmt"
	"os"
)

func isSorted(p []int) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != i+1 {
			return false
		}
	}
	return true
}

func op1(p []int) []int {
	res := make([]int, len(p))
	copy(res, p)
	for i := 0; i < len(p); i += 2 {
		res[i], res[i+1] = res[i+1], res[i]
	}
	return res
}

func op2(p []int) []int {
	res := make([]int, len(p))
	copy(res, p)
	n := len(p) / 2
	for i := 0; i < n; i++ {
		res[i], res[i+n] = res[i+n], res[i]
	}
	return res
}

func simulate(start []int, firstOp int) int {
	cur := make([]int, len(start))
	copy(cur, start)
	if isSorted(cur) {
		return 0
	}
	maxSteps := len(start) * 4
	op := firstOp
	for step := 1; step <= maxSteps; step++ {
		if op == 1 {
			cur = op1(cur)
		} else {
			cur = op2(cur)
		}
		if isSorted(cur) {
			return step
		}
		if op == 1 {
			op = 2
		} else {
			op = 1
		}
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, 2*n)
	for i := range p {
		fmt.Fscan(in, &p[i])
	}
	res1 := simulate(p, 1)
	res2 := simulate(p, 2)
	ans := -1
	if res1 != -1 {
		ans = res1
	}
	if res2 != -1 && (ans == -1 || res2 < ans) {
		ans = res2
	}
	fmt.Println(ans)
}
