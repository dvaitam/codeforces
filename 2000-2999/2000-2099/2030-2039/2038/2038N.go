package main

import (
	"bufio"
	"fmt"
	"os"
)

type expr struct {
	a, op, b byte
}

func cost(s expr, e expr) int {
	c := 0
	if s.a != e.a {
		c++
	}
	if s.op != e.op {
		c++
	}
	if s.b != e.b {
		c++
	}
	return c
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		best := expr{}
		bestCost := 4
		for a := byte('0'); a <= '9'; a++ {
			for opIdx := 0; opIdx < 3; opIdx++ {
				op := byte("<=>"[opIdx])
				for b := byte('0'); b <= '9'; b++ {
					valid := false
					switch op {
					case '<':
						valid = a < b
					case '=':
						valid = a == b
					case '>':
						valid = a > b
					}
					if !valid {
						continue
					}
					cur := expr{a, op, b}
					c := cost(expr{s[0], s[1], s[2]}, cur)
					if c < bestCost {
						bestCost = c
						best = cur
					}
				}
			}
		}
		res := []byte{best.a, best.op, best.b}
		fmt.Fprintln(out, string(res))
	}
}
