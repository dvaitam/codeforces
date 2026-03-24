package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// solve is an embedded reference solver for 1610H.
func solve(input string) string {
	r := strings.NewReader(input)
	var n, m int
	fmt.Fscan(r, &n, &m)

	adj := make([][]int, n+1)
	up := make([][20]int, n+1)
	depth := make([]int, n+1)

	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(r, &p)
		up[i][0] = p
		depth[i] = depth[p] + 1
		adj[p] = append(adj[p], i)
		for j := 1; j < 20; j++ {
			up[i][j] = up[up[i][j-1]][j-1]
		}
	}

	tin := make([]int, n+1)
	tout := make([]int, n+1)
	timer := 0

	var dfs func(v int)
	dfs = func(v int) {
		timer++
		tin[v] = timer
		for _, u := range adj[v] {
			dfs(u)
		}
		tout[v] = timer
	}
	dfs(1)

	isAncestor := func(u, v int) bool {
		return tin[u] <= tin[v] && tout[u] >= tout[v]
	}

	getChildAncestor := func(x, y int) int {
		curr := y
		for i := 19; i >= 0; i-- {
			if depth[curr]-(1<<i) > depth[x] {
				curr = up[curr][i]
			}
		}
		return curr
	}

	type Player struct {
		x, y, top, typ int
	}
	players := make([]Player, 0, m)

	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(r, &x, &y)

		if x == up[y][0] || y == up[x][0] {
			return "-1"
		}

		if isAncestor(y, x) {
			x, y = y, x
		}

		if isAncestor(x, y) {
			ux := getChildAncestor(x, y)
			players = append(players, Player{x, y, ux, 1})
		} else {
			players = append(players, Player{x, y, 1, 2})
		}
	}

	// sort by depth descending
	for i := 0; i < len(players); i++ {
		for j := i + 1; j < len(players); j++ {
			if depth[players[j].top] > depth[players[i].top] {
				players[i], players[j] = players[j], players[i]
			}
		}
	}

	bit := make([]int, n+2)
	add := func(idx, val int) {
		for ; idx <= n; idx += idx & -idx {
			bit[idx] += val
		}
	}
	query := func(idx int) int {
		sum := 0
		for ; idx > 0; idx -= idx & -idx {
			sum += bit[idx]
		}
		return sum
	}
	queryRange := func(l, rr int) int {
		if l > rr {
			return 0
		}
		return query(rr) - query(l-1)
	}

	ans := 0
	for _, p := range players {
		hit := false
		if p.typ == 1 {
			cntUx := queryRange(tin[p.top], tout[p.top])
			cntY := queryRange(tin[p.y], tout[p.y])
			if cntUx-cntY > 0 {
				hit = true
			}
		} else {
			cntT := queryRange(1, n)
			cntX := queryRange(tin[p.x], tout[p.x])
			cntY := queryRange(tin[p.y], tout[p.y])
			if cntT-cntX-cntY > 0 {
				hit = true
			}
		}
		if !hit {
			ans++
			add(tin[p.top], 1)
		}
	}

	return fmt.Sprint(ans)
}

func genTree(n int) []int {
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = rand.Intn(i-1) + 1
	}
	return parents
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(5) + 2
		m := rand.Intn(5) + 1
		parents := genTree(n)
		input := fmt.Sprintf("%d %d\n", n, m)
		for _, p := range parents {
			input += fmt.Sprintf("%d ", p)
		}
		input += "\n"
		for i := 0; i < m; i++ {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			input += fmt.Sprintf("%d %d\n", x, y)
		}
		exp := solve(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec failed: %v\n", t, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed expected %s got %s\n", t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
