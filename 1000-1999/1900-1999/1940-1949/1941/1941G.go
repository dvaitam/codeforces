package main

import (
	"bufio"
	"fmt"
	"os"
)

type deque struct {
	data       []int
	head, tail int
}

func newDeque(n int) *deque {
	size := 2*n + 10
	d := &deque{data: make([]int, size)}
	d.head = size / 2
	d.tail = d.head
	return d
}

func (d *deque) empty() bool { return d.head == d.tail }

func (d *deque) pushFront(x int) {
	d.head--
	d.data[d.head] = x
}

func (d *deque) pushBack(x int) {
	d.data[d.tail] = x
	d.tail++
}

func (d *deque) popFront() int {
	x := d.data[d.head]
	d.head++
	return x
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
		colorID := make(map[int]int)
		vertexColors := make([][]int, n+1)
		colorVerts := [][]int{{}}
		for i := 0; i < m; i++ {
			var u, v, c int
			fmt.Fscan(in, &u, &v, &c)
			id, ok := colorID[c]
			if !ok {
				id = len(colorVerts)
				colorID[c] = id
				colorVerts = append(colorVerts, []int{})
			}
			vertexColors[u] = append(vertexColors[u], id)
			vertexColors[v] = append(vertexColors[v], id)
			colorVerts[id] = append(colorVerts[id], u)
			colorVerts[id] = append(colorVerts[id], v)
		}
		var b, e int
		fmt.Fscan(in, &b, &e)
		if b == e {
			fmt.Fprintln(out, 0)
			continue
		}
		numColors := len(colorVerts) - 1
		total := n + numColors
		const INF = int(1e9)
		dist := make([]int, total+1)
		for i := 1; i <= total; i++ {
			dist[i] = INF
		}
		dq := newDeque(total)
		dist[b] = 0
		dq.pushFront(b)
		for !dq.empty() {
			v := dq.popFront()
			d := dist[v]
			if v == e {
				break
			}
			if v <= n {
				for _, cid := range vertexColors[v] {
					node := n + cid
					if dist[node] > d+1 {
						dist[node] = d + 1
						dq.pushBack(node)
					}
				}
			} else {
				cid := v - n
				if colorVerts[cid] != nil {
					for _, u := range colorVerts[cid] {
						if dist[u] > d {
							dist[u] = d
							dq.pushFront(u)
						}
					}
					colorVerts[cid] = nil
				}
			}
		}
		fmt.Fprintln(out, dist[e])
	}
}
