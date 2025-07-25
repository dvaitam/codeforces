package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1_000_000_007

type State struct {
	x, y int
	grid []string
}

func serialize(g []string) string {
	s := make([]byte, 0, len(g)*(len(g[0])+1))
	for _, row := range g {
		s = append(s, row...)
		s = append(s, '\n')
	}
	return string(s)
}

func countPaths(g []string) int {
	n := len(g)
	m := len(g[0])
	start := State{0, 0, append([]string(nil), g...)}
	q := []State{start}
	ways := map[string]int{serialize(start.grid) + "0,0": 1}
	res := 0

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		key := serialize(cur.grid) + fmt.Sprintf("%d,%d", cur.x, cur.y)
		w := ways[key]
		if cur.x == n-1 && cur.y == m-1 {
			res = (res + w) % MOD
			continue
		}
		for _, d := range [][2]int{{1, 0}, {0, 1}} {
			nx := cur.x + d[0]
			ny := cur.y + d[1]
			if nx >= n || ny >= m {
				continue
			}
			ng := make([]string, n)
			for i := range cur.grid {
				ng[i] = cur.grid[i]
			}
			if d[0] == 1 { // down
				if ng[nx][ny] == 'R' {
					k := nx
					for k < n && ng[k][ny] == 'R' {
						k++
					}
					if k == n {
						continue
					}
					b := []byte(ng[k])
					b[ny] = 'R'
					ng[k] = string(b)
					for t := k - 1; t >= nx; t-- {
						row := []byte(ng[t])
						row[ny] = '.'
						ng[t] = string(row)
					}
				}
			} else {
				if ng[nx][ny] == 'R' {
					k := ny
					for k < m && ng[nx][k] == 'R' {
						k++
					}
					if k == m {
						continue
					}
					row := []byte(ng[nx])
					row[k] = 'R'
					for t := k - 1; t >= ny; t-- {
						row[t] = '.'
					}
					ng[nx] = string(row)
				}
			}
			nk := serialize(ng) + fmt.Sprintf("%d,%d", nx, ny)
			if _, ok := ways[nk]; !ok {
				q = append(q, State{nx, ny, ng})
			}
			ways[nk] = (ways[nk] + w) % MOD
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	g := make([]string, n)
	for i := range g {
		fmt.Fscan(in, &g[i])
	}
	if g[0][0] == 'R' {
		fmt.Println(0)
		return
	}
	fmt.Println(countPaths(g))
}
