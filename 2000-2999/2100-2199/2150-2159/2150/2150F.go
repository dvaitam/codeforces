package main

import (
	"bufio"
	"fmt"
	"os"
)

type triple struct {
	a, b, c int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		e := make([][]bool, n+1)
		for i := range e {
			e[i] = make([]bool, n+1)
		}
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			e[u][v] = true
			e[v][u] = true
		}

		tr := make([][]int, n+1)
		vis := make([]bool, n+1)
		var dfs func(int)
		dfs = func(u int) {
			vis[u] = true
			for v := 1; v <= n; v++ {
				if !e[u][v] || vis[v] {
					continue
				}
				tr[u] = append(tr[u], v)
				tr[v] = append(tr[v], u)
				dfs(v)
			}
		}
		dfs(1)

		path := make([][][]int, n+1)
		for i := 0; i <= n; i++ {
			path[i] = make([][]int, n+1)
		}
		var dfs2 func(int, int, int, []int)
		dfs2 = func(u, parent, root int, cur []int) {
			cur = append(cur, u)
			path[root][u] = append([]int(nil), cur...)
			for _, v := range tr[u] {
				if v == parent {
					continue
				}
				dfs2(v, u, root, cur)
			}
		}
		for root := 1; root <= n; root++ {
			dfs2(root, 0, root, nil)
		}

		dia := -1
		dx, dy := 1, 1
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				d := len(path[i][j]) - 1
				if d > dia {
					dia = d
					dx, dy = i, j
				}
			}
		}
		lenVal := (dia + 1) / 2

		op1 := make([]triple, 0)
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				if e[i][j] {
					continue
				}
				for k := 1; k <= n; k++ {
					if e[i][k] && e[j][k] {
						op1 = append(op1, triple{i, k, j})
						break
					}
				}
			}
		}

		fmt.Fprintln(out, 2)
		fmt.Fprintln(out, 3)
		fmt.Fprintln(out, len(op1))
		for _, tp := range op1 {
			fmt.Fprintf(out, "%d %d %d\n", tp.a, tp.b, tp.c)
			e[tp.a][tp.c] = true
			e[tp.c][tp.a] = true
		}

		op2 := make([][]int, 0)
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				if e[i][j] {
					continue
				}
				p := constructPath(i, j, lenVal, path, dx, dy)
				op2 = append(op2, p)
			}
		}

		fmt.Fprintln(out, lenVal+1)
		fmt.Fprintln(out, len(op2))
		for _, p := range op2 {
			for idx, node := range p {
				if idx > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, node)
			}
			fmt.Fprintln(out)
		}
	}
}

func constructPath(u, v, lenVal int, path [][][]int, dx, dy int) []int {
	pv := path[u][v]
	dis := len(pv) - 1
	if dis >= lenVal {
		return buildLongPath(pv, lenVal)
	}
	nd := lenVal - dis
	xi, yi := dx, dy
	if len(path[u][xi]) < len(path[u][yi]) {
		xi, yi = yi, xi
	}
	pix := path[u][xi]
	t := 0
	for t+1 < len(pv) && t+1 < len(pix) && pv[t+1] == pix[t+1] {
		t++
	}
	res := make([]int, 0, lenVal+1)
	sameFirst := len(pix) <= 1 || len(pv) <= 1 || pix[1] == pv[1]
	if !sameFirst {
		res = append(res, u)
		for k := 2; k <= nd; k += 2 {
			res = append(res, pix[k])
		}
		start := nd
		if nd%2 == 0 {
			start--
		}
		for k := start; k >= 1; k -= 2 {
			res = append(res, pix[k])
		}
		for k := 1; k < len(pv); k++ {
			res = append(res, pv[k])
		}
		return res
	}
	for k := 0; k < t; k++ {
		res = append(res, pv[k])
	}
	suffix := pix[t:]
	for k := 1; k <= nd; k += 2 {
		res = append(res, suffix[k])
	}
	start := nd
	if nd%2 == 1 {
		start--
	}
	for k := start; k >= 1; k -= 2 {
		res = append(res, suffix[k])
	}
	for k := t; k < len(pv); k++ {
		res = append(res, pv[k])
	}
	return res
}

func buildLongPath(p []int, lenVal int) []int {
	dis := len(p) - 1
	res := make([]int, 0, lenVal+1)
	k := 0
	for {
		res = append(res, p[k])
		if len(res)+(dis-k) == lenVal+1 {
			for k < dis {
				k++
				res = append(res, p[k])
			}
			break
		}
		k += 2
	}
	return res
}
