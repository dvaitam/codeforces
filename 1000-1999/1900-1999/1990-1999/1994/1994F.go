package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	v, id int
}

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func nextInt() int {
	var b byte
	var err error
	// skip non-numeric
	for {
		b, err = reader.ReadByte()
		if err != nil {
			return 0
		}
		if (b >= '0' && b <= '9') || b == '-' {
			break
		}
	}
	sign := 1
	if b == '-' {
		sign = -1
		b, _ = reader.ReadByte()
	}
	x := 0
	for ; b >= '0' && b <= '9'; b, _ = reader.ReadByte() {
		x = x*10 + int(b-'0')
	}
	return x * sign
}

func main() {
	defer writer.Flush()
	t := nextInt()
	for tc := 0; tc < t; tc++ {
		n := nextInt()
		m := nextInt()
		finalG := make([][]edge, n)
		optionalG := make([][]edge, n)
		for i := 0; i < m; i++ {
			u := nextInt() - 1
			v := nextInt() - 1
			c := nextInt()
			if c == 1 {
				finalG[u] = append(finalG[u], edge{v, i})
				finalG[v] = append(finalG[v], edge{u, i})
			} else {
				optionalG[u] = append(optionalG[u], edge{v, i})
				optionalG[v] = append(optionalG[v], edge{u, i})
			}
		}
		trav := make([]bool, n)
		parent := make([]int, n)
		parentEdge := make([]int, n)
		iter := make([]int, n)
		odd := make([]int, n)
		for i := 0; i < n; i++ {
			odd[i] = len(finalG[i]) & 1
			parent[i] = -1
		}
		order := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if trav[i] {
				continue
			}
			trav[i] = true
			stack := []int{i}
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				if iter[u] < len(optionalG[u]) {
					e := optionalG[u][iter[u]]
					iter[u]++
					if !trav[e.v] {
						trav[e.v] = true
						parent[e.v] = u
						parentEdge[e.v] = e.id
						stack = append(stack, e.v)
					}
				} else {
					order = append(order, u)
					stack = stack[:len(stack)-1]
				}
			}
		}

		hasSol := true
		for _, u := range order {
			p := parent[u]
			if p == -1 {
				if odd[u] == 1 {
					hasSol = false
				}
				continue
			}
			if odd[u] == 1 {
				odd[p] ^= 1
				finalG[p] = append(finalG[p], edge{u, parentEdge[u]})
				finalG[u] = append(finalG[u], edge{p, parentEdge[u]})
			}
		}
		if !hasSol {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		used := make([]bool, m)
		ptr := make([]int, n)
		var ans []int
		stack := []int{0}
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			for ptr[u] < len(finalG[u]) && used[finalG[u][ptr[u]].id] {
				ptr[u]++
			}
			if ptr[u] == len(finalG[u]) {
				ans = append(ans, u+1)
				stack = stack[:len(stack)-1]
				continue
			}
			e := finalG[u][ptr[u]]
			ptr[u]++
			used[e.id] = true
			stack = append(stack, e.v)
		}
		// output edges count = len(ans)-1
		fmt.Fprintln(writer, len(ans)-1)
		for i, v := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		writer.WriteByte('\n')
	}
}
