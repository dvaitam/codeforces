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

type dsu struct {
	p []int
}

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) {
	fa, fb := d.find(a), d.find(b)
	if fa != fb {
		d.p[fa] = fb
	}
}

func expected(n int, edges [][3]int, queries [][2]int) []int {
	// colors from 1..m maybe
	maxColor := 0
	for _, e := range edges {
		if e[2] > maxColor {
			maxColor = e[2]
		}
	}
	comps := make([]*dsu, maxColor+1)
	for c := 1; c <= maxColor; c++ {
		comps[c] = newDSU(n)
	}
	for _, e := range edges {
		comps[e[2]].union(e[0], e[1])
	}
	res := make([]int, len(queries))
	for i, q := range queries {
		cnt := 0
		for c := 1; c <= maxColor; c++ {
			if comps[c].find(q[0]) == comps[c].find(q[1]) {
				cnt++
			}
		}
		res[i] = cnt
	}
	return res
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(5) + 2
		m := rng.Intn(6) + 1
		used := map[[3]int]struct{}{}
		edges := make([][3]int, 0, m)
		for len(edges) < m {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			if a == b {
				continue
			}
			c := rng.Intn(m) + 1
			key := [3]int{a, b, c}
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			edges = append(edges, [3]int{a, b, c})
		}
		q := rng.Intn(5) + 1
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for u == v {
				v = rng.Intn(n) + 1
			}
			queries[i] = [2]int{u, v}
		}
		// build input string
		input := fmt.Sprintf("%d %d\n", n, m)
		for _, e := range edges {
			input += fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2])
		}
		input += fmt.Sprintf("%d\n", q)
		for _, qu := range queries {
			input += fmt.Sprintf("%d %d\n", qu[0], qu[1])
		}
		expect := expected(n, edges, queries)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		parts := strings.Fields(out)
		if len(parts) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", tc+1, len(expect), len(parts), input)
			os.Exit(1)
		}
		for i, exp := range expect {
			if parts[i] != fmt.Sprintf("%d", exp) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", tc+1, exp, parts[i], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
