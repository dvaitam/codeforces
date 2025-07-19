package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
	c, _ := reader.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = reader.ReadByte()
	}
	sign := 1
	if c == '-' {
		sign = -1
		c, _ = reader.ReadByte()
	}
	x := 0
	for c >= '0' && c <= '9' {
		x = x*10 + int(c-'0')
		c, _ = reader.ReadByte()
	}
	return x * sign
}

func main() {
	defer writer.Flush()
	n := readInt()
	q := readInt()
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		fa := readInt()
		children[fa] = append(children[fa], i)
	}
	for i := 1; i <= n; i++ {
		if len(children[i]) > 1 {
			sort.Ints(children[i])
		}
	}
	st := make([]int, n+1)
	ed := make([]int, n+1)
	id := make([]int, n+2)
	num := 0
	type frame struct{ node, idx int }
	stack := make([]frame, 0, n)
	stack = append(stack, frame{1, 0})
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.idx == 0 {
			num++
			st[top.node] = num
			id[num] = top.node
		}
		if top.idx < len(children[top.node]) {
			child := children[top.node][top.idx]
			top.idx++
			stack = append(stack, frame{child, 0})
		} else {
			ed[top.node] = num
			stack = stack[:len(stack)-1]
		}
	}
	for i := 0; i < q; i++ {
		u := readInt()
		v := readInt()
		pos := st[u] + v - 1
		ans := -1
		if pos <= ed[u] {
			ans = id[pos]
		}
		writer.WriteString(strconv.Itoa(ans))
		writer.WriteByte('\n')
	}
}
