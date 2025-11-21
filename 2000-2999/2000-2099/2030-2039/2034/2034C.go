package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastQueue struct {
	data []int
	head int
}

func (q *FastQueue) Push(x int) {
	q.data = append(q.data, x)
}

func (q *FastQueue) Pop() int {
	x := q.data[q.head]
	q.head++
	if q.head*2 >= len(q.data) {
		q.data = q.data[q.head:]
		q.head = 0
	}
	return x
}

func (q *FastQueue) Empty() bool {
	return q.head >= len(q.data)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}
		total := n * m
		fixedPred := make([][]int, total)
		deps := make([][]int, total)
		fixedTarget := make([]int, total)
		for i := range fixedTarget {
			fixedTarget[i] = -2
		}
		qCount := make([]int, total)
		removed := make([]bool, total)
		q := FastQueue{data: make([]int, 0)}

		idx := func(r, c int) int {
			return r*m + c
		}

		dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

		for r := 0; r < n; r++ {
			for c := 0; c < m; c++ {
				id := idx(r, c)
				ch := grid[r][c]
				if ch == '?' {
					neighbors := 0
					for _, d := range dirs {
						nr, nc := r+d[0], c+d[1]
						if 0 <= nr && nr < n && 0 <= nc && nc < m {
							neighbors++
							v := idx(nr, nc)
							deps[v] = append(deps[v], id)
						}
					}
					qCount[id] = neighbors
					if neighbors == 0 {
						q.Push(id)
					}
				} else {
					nr, nc := r, c
					switch ch {
					case 'U':
						nr--
					case 'D':
						nr++
					case 'L':
						nc--
					case 'R':
						nc++
					}
					if nr < 0 || nr >= n || nc < 0 || nc >= m {
						fixedTarget[id] = -1
						q.Push(id)
					} else {
						v := idx(nr, nc)
						fixedTarget[id] = v
						fixedPred[v] = append(fixedPred[v], id)
					}
				}
			}
		}

		for !q.Empty() {
			u := q.Pop()
			if removed[u] {
				continue
			}
			removed[u] = true
			for _, pred := range fixedPred[u] {
				if !removed[pred] {
					q.Push(pred)
				}
			}
			for _, dep := range deps[u] {
				if removed[dep] {
					continue
				}
				qCount[dep]--
				if qCount[dep] == 0 {
					q.Push(dep)
				}
			}
		}

		trapped := 0
		for _, rem := range removed {
			if !rem {
				trapped++
			}
		}
		fmt.Fprintln(out, trapped)
	}
}
