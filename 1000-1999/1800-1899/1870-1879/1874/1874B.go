package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct{ x, y uint32 }

type queue struct{ data []state }

func (q *queue) push(s state) { q.data = append(q.data, s) }
func (q *queue) pop() state   { v := q.data[0]; q.data = q.data[1:]; return v }

func solve(a, b, c, d, m uint32) int {
	start := state{a, b}
	target := state{c, d}
	if start == target {
		return 0
	}
	q := queue{data: []state{start}}
	dist := map[state]int{start: 0}
	for len(q.data) > 0 {
		cur := q.pop()
		step := dist[cur]
		if cur == target {
			return step
		}
		x, y := cur.x, cur.y
		next := [4]state{
			{x & y, y},
			{x | y, y},
			{x, x ^ y},
			{x, y ^ m},
		}
		for _, v := range next {
			if _, ok := dist[v]; !ok {
				dist[v] = step + 1
				if v == target {
					return step + 1
				}
				q.push(v)
			}
		}
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b, c, d, m uint32
		fmt.Fscan(reader, &a, &b, &c, &d, &m)
		fmt.Fprintln(writer, solve(a, b, c, d, m))
	}
}
