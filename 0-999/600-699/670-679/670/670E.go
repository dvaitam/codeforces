package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, p int
	if _, err := fmt.Fscan(in, &n, &m, &p); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	var ops string
	fmt.Fscan(in, &ops)

	pair := make([]int, n+1)
	stack := make([]int, 0, n/2)
	for i := 1; i <= n; i++ {
		if s[i-1] == '(' {
			stack = append(stack, i)
		} else {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pair[i] = j
			pair[j] = i
		}
	}

	next := make([]int, n+2)
	prev := make([]int, n+2)
	for i := 1; i <= n; i++ {
		next[i] = i + 1
		prev[i] = i - 1
	}
	next[n] = 0
	prev[1] = 0

	cur := p
	for i := 0; i < len(ops); i++ {
		switch ops[i] {
		case 'L':
			cur = prev[cur]
		case 'R':
			cur = next[cur]
		case 'D':
			l := cur
			r := pair[cur]
			if l > r {
				l, r = r, l
			}
			lp := prev[l]
			rn := next[r]
			if lp != 0 {
				next[lp] = rn
			}
			if rn != 0 {
				prev[rn] = lp
			}
			if rn != 0 {
				cur = rn
			} else {
				cur = lp
			}
		}
	}

	for prev[cur] != 0 {
		cur = prev[cur]
	}

	out := bufio.NewWriter(os.Stdout)
	for i := cur; i != 0; i = next[i] {
		out.WriteByte(s[i-1])
	}
	out.WriteByte('\n')
	out.Flush()
}
