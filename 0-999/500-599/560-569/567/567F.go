package main

import (
	"bufio"
	"fmt"
	"os"
)

type Constraint struct {
	x  int
	op string
	y  int
}

var (
	n     int
	k     int
	cons  []Constraint
	count int64
	a     []int
)

func check(seq []int) bool {
	for _, c := range cons {
		vx := seq[c.x-1]
		vy := seq[c.y-1]
		switch c.op {
		case "=":
			if vx != vy {
				return false
			}
		case "<":
			if !(vx < vy) {
				return false
			}
		case ">":
			if !(vx > vy) {
				return false
			}
		case "<=":
			if !(vx <= vy) {
				return false
			}
		case ">=":
			if !(vx >= vy) {
				return false
			}
		}
	}
	return true
}

func dfs(i int) {
	if i > n {
		// build sequence
		left := make([]int, 0, 2*n)
		for v := 1; v <= n; v++ {
			for j := 0; j < a[v]; j++ {
				left = append(left, v)
			}
		}
		right := make([]int, 0, 2*n)
		for v := n; v >= 1; v-- {
			for j := 0; j < 2-a[v]; j++ {
				right = append(right, v)
			}
		}
		seq := append(left, right...)
		if check(seq) {
			count++
		}
		return
	}
	for t := 0; t <= 2; t++ {
		a[i] = t
		dfs(i + 1)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &k)
	cons = make([]Constraint, k)
	for i := 0; i < k; i++ {
		var x, y int
		var op string
		fmt.Fscan(in, &x, &op, &y)
		cons[i] = Constraint{x: x, op: op, y: y}
	}
	a = make([]int, n+1)
	dfs(1)
	fmt.Println(count)
}
