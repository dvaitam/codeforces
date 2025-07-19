package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	// fast int reader
	readInt := func() int {
		var sb strings.Builder
		for {
			b, err := reader.ReadByte()
			if err != nil {
				break
			}
			if (b >= '0' && b <= '9') || b == '-' {
				sb.WriteByte(b)
				break
			}
		}
		for {
			b, err := reader.ReadByte()
			if err != nil || b < '0' || b > '9' {
				break
			}
			sb.WriteByte(b)
		}
		x, _ := strconv.Atoi(sb.String())
		return x
	}
	n := readInt()
	m := readInt()
	k := readInt()
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		u := readInt()
		v := readInt()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// DFS tree iterative
	vis := make([]bool, n+1)
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	sonCount := make([]int, n+1)
	var leaves []int
	type stkItem struct{ u, idx int }
	stack := make([]stkItem, 0, n)
	vis[1] = true
	depth[1] = 1
	stack = append(stack, stkItem{1, 0})
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u := top.u
		if top.idx < len(adj[u]) {
			v := adj[u][top.idx]
			top.idx++
			if !vis[v] {
				vis[v] = true
				parent[v] = u
				depth[v] = depth[u] + 1
				sonCount[u]++
				stack = append(stack, stkItem{v, 0})
			}
		} else {
			if sonCount[u] == 0 {
				leaves = append(leaves, u)
			}
			stack = stack[:len(stack)-1]
		}
	}
	// find farthest
	o := 1
	for i := 1; i <= n; i++ {
		if depth[i] > depth[o] {
			o = i
		}
	}
	// check path
	if depth[o]*k >= n {
		// PATH
		writer.WriteString("PATH\n")
		var path []int
		for u := o; u != 0; u = parent[u] {
			path = append(path, u)
		}
		// output
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(len(path)))
		sb.WriteByte('\n')
		for _, v := range path {
			sb.WriteString(strconv.Itoa(v))
			sb.WriteByte(' ')
		}
		sb.WriteByte('\n')
		writer.WriteString(sb.String())
	} else {
		// CYCLES
		writer.WriteString("CYCLES\n")
		// create k cycles
		for i := 0; i < k; i++ {
			u := leaves[i]
			var T []int
			for _, g := range adj[u] {
				if parent[g] != u && g != parent[u] {
					T = append(T, g)
					if len(T) >= 2 {
						break
					}
				}
			}
			a, b := T[0], T[1]
			if depth[a] > depth[b] {
				a, b = b, a
			}
			var cycle []int
			if (depth[u]-depth[a]+1)%3 != 0 {
				for v := u; v != a; v = parent[v] {
					cycle = append(cycle, v)
				}
				cycle = append(cycle, a)
			} else if (depth[u]-depth[b]+1)%3 != 0 {
				for v := u; v != b; v = parent[v] {
					cycle = append(cycle, v)
				}
				cycle = append(cycle, b)
			} else {
				cycle = append(cycle, u)
				for v := b; v != a; v = parent[v] {
					cycle = append(cycle, v)
				}
				cycle = append(cycle, a)
			}
			// output each cycle
			var sb strings.Builder
			sb.WriteString(strconv.Itoa(len(cycle)))
			sb.WriteByte('\n')
			for _, v := range cycle {
				sb.WriteString(strconv.Itoa(v))
				sb.WriteByte(' ')
			}
			sb.WriteByte('\n')
			writer.WriteString(sb.String())
		}
	}
}
