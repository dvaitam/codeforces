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

type edge struct{ u, v int }

func expected(n int, p []int, edges []edge) string {
	adj := make(map[int]map[int]bool)
	for _, e := range edges {
		if adj[e.u] == nil {
			adj[e.u] = make(map[int]bool)
		}
		adj[e.u][e.v] = true
	}
	bj := make(map[int]bool)
	last := p[n-1]
	for _, e := range edges {
		if e.v == last {
			bj[e.u] = true
		}
	}
	var a []int
	ans := 0
	for i := n - 2; i >= 0; i-- {
		pi := p[i]
		if bj[pi] {
			ans++
			ok := true
			for _, idx := range a {
				if !adj[pi][p[idx]] {
					ok = false
					break
				}
			}
			if !ok {
				ans--
				a = append(a, i)
			}
		} else {
			a = append(a, i)
		}
	}
	return fmt.Sprintf("%d", ans)
}

func runCase(bin, input, want string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n     int
		p     []int
		edges []edge
	}
	tests := []test{
		{n: 2, p: []int{1, 2}, edges: []edge{{1, 2}}},
		{n: 3, p: []int{1, 2, 3}, edges: []edge{{2, 3}}},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		perm := rand.Perm(n)
		for i := range perm {
			perm[i]++
		}
		m := rng.Intn(n*n + 1)
		es := make([]edge, m)
		for j := 0; j < m; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				if u < n {
					v = u + 1
				} else {
					v = u - 1
				}
			}
			es[j] = edge{u, v}
		}
		tests = append(tests, test{n: n, p: perm, edges: es})
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		want := expected(tc.n, tc.p, tc.edges)
		if err := runCase(bin, sb.String(), want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
