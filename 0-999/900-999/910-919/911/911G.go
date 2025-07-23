package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxVal = 101

// Block represents a chunk of the array with lazy mapping using DSU
type Block struct {
	l, r   int
	a      []int
	parent [maxVal]int
}

func (b *Block) initParent() {
	for i := 1; i < maxVal; i++ {
		b.parent[i] = i
	}
}

func (b *Block) find(x int) int {
	if b.parent[x] != x {
		b.parent[x] = b.find(b.parent[x])
	}
	return b.parent[x]
}

func (b *Block) union(x, y int) {
	px := b.find(x)
	py := b.find(y)
	if px != py {
		b.parent[px] = py
	}
}

func (b *Block) push() {
	for i := range b.a {
		b.a[i] = b.find(b.a[i])
	}
	b.initParent()
}

func (b *Block) update(l, r, x, y int) {
	if l <= b.l && b.r <= r {
		b.union(x, y)
		return
	}
	b.push()
	if l < b.l {
		l = b.l
	}
	if r > b.r {
		r = b.r
	}
	for i := l; i <= r; i++ {
		idx := i - b.l
		if b.a[idx] == x {
			b.a[idx] = y
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	var q int
	fmt.Fscan(in, &q)

	blockSize := 450
	blocks := make([]*Block, 0)
	for start := 0; start < n; {
		end := start + blockSize
		if end > n {
			end = n
		}
		b := &Block{l: start, r: end - 1, a: append([]int(nil), arr[start:end]...)}
		b.initParent()
		blocks = append(blocks, b)
		start = end
	}

	for ; q > 0; q-- {
		var l, r, x, y int
		fmt.Fscan(in, &l, &r, &x, &y)
		if x == y {
			continue
		}
		l--
		r--
		for _, b := range blocks {
			if b.r < l || b.l > r {
				continue
			}
			b.update(l, r, x, y)
		}
	}

	for idx, b := range blocks {
		b.push()
		for i, v := range b.a {
			if idx > 0 || i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
	}
	fmt.Fprintln(out)
}
