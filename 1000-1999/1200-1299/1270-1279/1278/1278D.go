package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int) {
	fx := d.find(x)
	fy := d.find(y)
	if fx == fy {
		return
	}
	if d.size[fx] < d.size[fy] {
		fx, fy = fy, fx
	}
	d.parent[fy] = fx
	d.size[fx] += d.size[fy]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	type Event struct {
		id   int
		left bool
	}
	events := make([]Event, 2*n+2) // 1..2n inclusive

	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		events[l] = Event{id: i + 1, left: true}
		events[r] = Event{id: i + 1, left: false}
	}

	dsu := NewDSU(n)
	stack := make([]int, 0)
	edges := 0

	for pos := 1; pos <= 2*n; pos++ {
		e := events[pos]
		if e.id == 0 {
			continue
		}
		id := e.id - 1
		if e.left {
			stack = append(stack, id)
		} else {
			temp := []int{}
			for len(stack) > 0 && stack[len(stack)-1] != id {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				dsu.union(id, top)
				edges++
				temp = append(temp, top)
				if edges >= n {
					fmt.Fprintln(writer, "NO")
					return
				}
			}
			if len(stack) == 0 {
				fmt.Fprintln(writer, "NO")
				return
			}
			// remove id
			stack = stack[:len(stack)-1]
			for i := len(temp) - 1; i >= 0; i-- {
				stack = append(stack, temp[i])
			}
		}
	}

	if edges != n-1 {
		fmt.Fprintln(writer, "NO")
		return
	}
	root := dsu.find(0)
	for i := 1; i < n; i++ {
		if dsu.find(i) != root {
			fmt.Fprintln(writer, "NO")
			return
		}
	}
	fmt.Fprintln(writer, "YES")
}
