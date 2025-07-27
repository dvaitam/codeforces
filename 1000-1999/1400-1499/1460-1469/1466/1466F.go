package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

type DSU struct {
	parent []int
}

func newDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &DSU{parent: p}
}

func (d *DSU) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) union(a, b int) bool {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return false
	}
	d.parent[b] = a
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	d := newDSU(m + 1) // include node 0
	used := make([]int, 0)
	for i := 1; i <= n; i++ {
		var k int
		fmt.Fscan(reader, &k)
		if k == 1 {
			var x int
			fmt.Fscan(reader, &x)
			if d.union(0, x) {
				used = append(used, i)
			}
		} else if k == 2 {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if d.union(x, y) {
				used = append(used, i)
			}
		}
	}

	rank := len(used)
	pow := int64(1)
	for i := 0; i < rank; i++ {
		pow = (pow * 2) % mod
	}

	fmt.Fprintf(writer, "%d %d\n", pow, rank)
	for i, v := range used {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}
