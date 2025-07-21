package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	first  int
	second int
}

func ask(w *bufio.Writer, r *bufio.Reader, x int) int {
	fmt.Fprintf(w, "? %d\n", x+1)
	w.Flush()
	var b int
	fmt.Fscan(r, &b)
	return b
}

func solve(g []map[int]struct{}, r *bufio.Reader, w *bufio.Writer) int {
	n := len(g)
	f := make([]int, n)
	d := make([]int, n)
	m := make([]int, n)
	var dfs func(int, int)
	dfs = func(u, p int) {
		f[u] = p
		m[u] = d[u]
		for v := range g[u] {
			if v != p {
				d[v] = d[u] + 1
				dfs(v, u)
				if m[v] > m[u] {
					m[u] = m[v]
				}
			}
		}
	}
	dfs(0, 0)
	b := make([]bool, n)
	l := 0
	for i := 1; i < n; i++ {
		if d[i] > d[l] {
			l = i
		}
	}
	if ask(w, r, l) != 0 {
		return l
	}
	for i := 0; i < n; i++ {
		if m[i]-d[i] == 51 && !b[i] {
			b[i] = true
			if ask(w, r, i) != 0 {
				if i != 0 {
					for {
						ask(w, r, l)
						if ask(w, r, i) == 0 {
							break
						}
					}
					if ask(w, r, f[i]) != 0 {
						return f[i]
					}
					return f[f[f[i]]]
				}
				break
			} else {
				delete(g[f[i]], i)
				delete(g[i], f[i])
				dfs(0, 0)
				i = -1
			}
		}
	}
	for i := 0; i < 100; i++ {
		ask(w, r, l)
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		g := make([]map[int]struct{}, n)
		for i := 0; i < n; i++ {
			g[i] = make(map[int]struct{})
		}
		for i := 1; i < n; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u][v] = struct{}{}
			g[v][u] = struct{}{}
		}
		x := solve(g, reader, writer)
		fmt.Fprintf(writer, "! %d\n", x+1)
	}
}
