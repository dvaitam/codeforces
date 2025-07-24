package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{make([]int, n+2)}
}

func (f *Fenwick) Add(i, delta int) {
	for i < len(f.tree) {
		f.tree[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.tree[i]
		i -= i & -i
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &a[i])
	}

	size := n + m + 2
	bit := NewFenwick(size)

	pos := make([]int, n+1)
	mn := make([]int, n+1)
	mx := make([]int, n+1)

	for i := 1; i <= n; i++ {
		pos[i] = m + i
		bit.Add(pos[i], 1)
		mn[i] = i
		mx[i] = i
	}

	cur := m
	for _, x := range a {
		// position before move
		currentPos := bit.Sum(pos[x])
		if mx[x] < currentPos {
			mx[x] = currentPos
		}
		mn[x] = 1
		bit.Add(pos[x], -1)
		pos[x] = cur
		bit.Add(pos[x], 1)
		cur--
	}

	for i := 1; i <= n; i++ {
		finalPos := bit.Sum(pos[i])
		if mx[i] < finalPos {
			mx[i] = finalPos
		}
		fmt.Fprintln(writer, mn[i], mx[i])
	}
}
