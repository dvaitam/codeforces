package main

import (
	"bufio"
	"os"
	"strconv"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
	x := 0
	sign := 1
	b, _ := reader.ReadByte()
	for (b < '0' || b > '9') && b != '-' {
		b, _ = reader.ReadByte()
	}
	if b == '-' {
		sign = -1
		b, _ = reader.ReadByte()
	}
	for b >= '0' && b <= '9' {
		x = x*10 + int(b-'0')
		b, _ = reader.ReadByte()
	}
	return x * sign
}

func main() {
	defer writer.Flush()
	T := readInt()
	for t := 0; t < T; t++ {
		n := readInt()
		m := readInt()
		g := make([][]int, n+1)
		for i := 0; i < m; i++ {
			u := readInt()
			v := readInt()
			if u != v {
				g[u] = append(g[u], v)
			}
		}
		dfn := make([]int, n+1)
		low := make([]int, n+1)
		col := make([]int, n+1)
		ins := make([]bool, n+1)
		stack := make([]int, 0, n)
		time := 0
		comp := 0
		var tarjan func(u int)
		tarjan = func(u int) {
			time++
			dfn[u] = time
			low[u] = time
			ins[u] = true
			stack = append(stack, u)
			for _, v := range g[u] {
				if dfn[v] == 0 {
					tarjan(v)
					if low[v] < low[u] {
						low[u] = low[v]
					}
				} else if ins[v] && dfn[v] < low[u] {
					low[u] = dfn[v]
				}
			}
			if low[u] == dfn[u] {
				comp++
				for {
					w := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					ins[w] = false
					col[w] = comp
					if w == u {
						break
					}
				}
			}
		}
		for i := 1; i <= n; i++ {
			if dfn[i] == 0 {
				tarjan(i)
			}
		}
		if comp == 1 {
			writer.WriteString("No\n")
			continue
		}
		writer.WriteString("Yes\n")
		outdeg := make([]int, comp+1)
		for u := 1; u <= n; u++ {
			for _, v := range g[u] {
				if col[u] != col[v] {
					outdeg[col[u]]++
				}
			}
		}
		id := 1
		for i := 1; i <= comp; i++ {
			if outdeg[i] == 0 {
				id = i
				break
			}
		}
		group1 := make([]int, 0)
		group2 := make([]int, 0)
		for i := 1; i <= n; i++ {
			if col[i] == id {
				group1 = append(group1, i)
			} else {
				group2 = append(group2, i)
			}
		}
		// print sizes
		writer.WriteString(strconv.Itoa(len(group1)))
		writer.WriteByte(' ')
		writer.WriteString(strconv.Itoa(len(group2)))
		writer.WriteByte('\n')
		// print group1
		for i, v := range group1 {
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(strconv.Itoa(v))
		}
		writer.WriteByte('\n')
		// print group2
		for i, v := range group2 {
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(strconv.Itoa(v))
		}
		writer.WriteByte('\n')
	}
}
