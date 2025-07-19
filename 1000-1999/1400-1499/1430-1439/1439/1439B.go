package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	check := make([]bool, n+1)
	used := make([]bool, n+1)
	edge := make(map[uint64]struct{}, m*2)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
		deg[a]++
		deg[b]++
		if a > b {
			a, b = b, a
		}
		key := uint64(a)<<32 | uint64(b)
		edge[key] = struct{}{}
	}
	// quick check
	if (k-5 > 0 && float64(k-5) > 2*math.Sqrt(float64(m))) || k-5 > n {
		fmt.Fprintln(writer, -1)
		return
	}
	// k-core reduction
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if deg[i] < k {
			queue = append(queue, i)
		}
	}
	qi := 0
	for qi < len(queue) {
		x := queue[qi]
		qi++
		if check[x] {
			continue
		}
		check[x] = true
		if deg[x] < k-1 {
			for _, y := range adj[x] {
				if !check[y] {
					deg[y]--
					if deg[y] == k-1 {
						queue = append(queue, y)
					}
				}
			}
			continue
		}
		if deg[x] == k-1 {
			save := make([]int, 0, k)
			save = append(save, x)
			used[x] = true
			for _, y := range adj[x] {
				if !check[y] {
					used[y] = true
					deg[y]--
					if deg[y] == k-1 {
						queue = append(queue, y)
					}
					save = append(save, y)
				}
			}
			f := false
			// check clique
			for i := 1; i < len(save); i++ {
				if k < 100 {
					for j := 0; j < i; j++ {
						a, b := save[i], save[j]
						if a > b {
							a, b = b, a
						}
						key := uint64(a)<<32 | uint64(b)
						if _, ok := edge[key]; !ok {
							f = true
							break
						}
					}
					if f {
						break
					}
				} else {
					cnt := 0
					for _, y := range adj[save[i]] {
						if used[y] {
							cnt++
						}
					}
					if cnt < k-1 {
						f = true
						break
					}
				}
			}
			// cleanup used
			for _, v := range save {
				used[v] = false
			}
			if !f {
				// found clique
				writer.WriteString("2\n")
				for idx, v := range save {
					if idx > 0 {
						writer.WriteString(" ")
					}
					writer.WriteString(fmt.Sprint(v))
				}
				writer.WriteString("\n")
				return
			}
		}
	}
	// remaining k-core
	remaining := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if !check[i] {
			remaining = append(remaining, i)
		}
	}
	if len(remaining) == 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	// output k-core
	writer.WriteString("1 ")
	writer.WriteString(fmt.Sprint(len(remaining)))
	writer.WriteString("\n")
	for i, v := range remaining {
		if i > 0 {
			writer.WriteString(" ")
		}
		writer.WriteString(fmt.Sprint(v))
	}
	writer.WriteString("\n")
}
