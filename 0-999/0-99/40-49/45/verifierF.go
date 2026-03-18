package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Embedded correct solver for 45F.
func solveF(input string) string {
	var m, n int
	fmt.Sscanf(input, "%d %d", &m, &n)

	minFunc := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	maxFunc := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	min3 := func(a, b, c int) int {
		return minFunc(a, minFunc(b, c))
	}
	ceilDiv2 := func(a int) int {
		if a <= 0 {
			return 0
		}
		return (a + 1) / 2
	}

	var parent [3][2][100005]int
	var vis [3][2][100005]bool

	for t := 0; t < 3; t++ {
		for b := 0; b < 2; b++ {
			for k := 0; k <= m+1; k++ {
				parent[t][b][k] = k
			}
		}
	}

	var find func(t, b, k int) int
	find = func(t, b, k int) int {
		if parent[t][b][k] == k {
			return k
		}
		parent[t][b][k] = find(t, b, parent[t][b][k])
		return parent[t][b][k]
	}

	type TK struct{ t, k int }

	getTypes := func(G, W int) []TK {
		var res []TK
		if G == 0 {
			res = append(res, TK{0, W})
		}
		if G == m {
			res = append(res, TK{2, W})
		}
		if G == W {
			res = append(res, TK{1, W})
		}
		return res
	}

	type Range struct{ t, L, R int }

	getB0Transitions := func(t, k int) []Range {
		var res []Range
		if t == 0 {
			if k >= 1 {
				res = append(res, Range{0, maxFunc(0, k-n), k - 1})
			}
		} else if t == 2 {
			if k >= 1 {
				res = append(res, Range{2, maxFunc(0, k-n), k - 1})
			}
			L := maxFunc(0, ceilDiv2(m+k-n))
			R := minFunc(k, (m+k-1)/2)
			if L <= R {
				res = append(res, Range{1, L, R})
			}
			if n >= m {
				maxY := min3(k, n-m, m)
				if maxY >= 0 {
					res = append(res, Range{0, k - maxY, k})
				}
			}
		} else if t == 1 {
			if k >= 1 && n >= 2 {
				res = append(res, Range{1, maxFunc(0, k-n/2), k - 1})
			}
			if k >= 1 && n >= k {
				res = append(res, Range{0, maxFunc(0, k-(n-k)), k})
			}
		}
		return res
	}

	getGW := func(t, k int) (int, int) {
		if t == 0 {
			return 0, k
		}
		if t == 1 {
			return k, k
		}
		return m, k
	}

	markVisited := func(G, W, b int) {
		for _, tk := range getTypes(G, W) {
			if !vis[tk.t][b][tk.k] {
				vis[tk.t][b][tk.k] = true
				parent[tk.t][b][tk.k] = find(tk.t, b, tk.k+1)
			}
		}
	}

	type State struct {
		G, W, b, dist int
	}

	queue := make([]State, 0, 1000000)
	queue = append(queue, State{m, m, 0, 0})
	markVisited(m, m, 0)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.G == 0 && curr.W == 0 && curr.b == 1 {
			return fmt.Sprintf("%d", curr.dist)
		}

		var targets []Range
		if curr.b == 0 {
			for _, tk := range getTypes(curr.G, curr.W) {
				targets = append(targets, getB0Transitions(tk.t, tk.k)...)
			}
		} else {
			Gd := m - curr.G
			Wd := m - curr.W
			for _, tkd := range getTypes(Gd, Wd) {
				targetsd := getB0Transitions(tkd.t, tkd.k)
				for _, r := range targetsd {
					var st int
					if r.t == 0 {
						st = 2
					} else if r.t == 1 {
						st = 1
					} else if r.t == 2 {
						st = 0
					}
					targets = append(targets, Range{st, m - r.R, m - r.L})
				}
			}
		}

		for _, r := range targets {
			k := find(r.t, 1-curr.b, r.L)
			for k <= r.R {
				Gnew, Wnew := getGW(r.t, k)
				queue = append(queue, State{Gnew, Wnew, 1 - curr.b, curr.dist + 1})
				markVisited(Gnew, Wnew, 1-curr.b)
				k = find(r.t, 1-curr.b, k+1)
			}
		}
	}

	return "-1"
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	m := rng.Intn(4) + 1
	n := rng.Intn(5) + 2
	return fmt.Sprintf("%d %d\n", m, n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solveF(in)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
