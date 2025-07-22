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

type Query struct {
	op int
	x  int
	y  int
}

type Test struct{ queries []Query }

func hasPath(l, r []int, a, b int) bool {
	n := len(l) - 1
	vis := make([]bool, n+1)
	q := []int{a}
	vis[a] = true
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if v == b {
			return true
		}
		for i := 1; i <= n; i++ {
			if !vis[i] && ((l[i] < l[v] && l[v] < r[i]) || (l[i] < r[v] && r[v] < r[i])) {
				vis[i] = true
				q = append(q, i)
			}
		}
	}
	return false
}

func expected(t Test) []string {
	l := []int{0}
	r := []int{0}
	res := []string{}
	for _, q := range t.queries {
		if q.op == 1 {
			l = append(l, q.x)
			r = append(r, q.y)
		} else {
			if hasPath(l, r, q.x, q.y) {
				res = append(res, "YES")
			} else {
				res = append(res, "NO")
			}
		}
	}
	return res
}

func runCase(bin string, t Test) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(t.queries)))
	for _, q := range t.queries {
		if q.op == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d\n", q.x, q.y))
		} else {
			sb.WriteString(fmt.Sprintf("2 %d %d\n", q.x, q.y))
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.Fields(strings.TrimSpace(out.String()))
	exp := expected(t)
	if len(got) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(got))
	}
	for i := range exp {
		if got[i] != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], got[i])
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) Test {
	m := rng.Intn(15) + 1
	t := Test{queries: make([]Query, 0, m)}
	cnt := 0
	maxLen := 0
	for i := 0; i < m; i++ {
		if cnt == 0 || rng.Intn(2) == 0 {
			length := maxLen + rng.Intn(5) + 1
			x := rng.Intn(50)
			y := x + length
			t.queries = append(t.queries, Query{1, x, y})
			cnt++
			maxLen = length
		} else {
			a := rng.Intn(cnt) + 1
			b := rng.Intn(cnt) + 1
			for b == a {
				b = rng.Intn(cnt) + 1
			}
			t.queries = append(t.queries, Query{2, a, b})
		}
	}
	return t
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []Test{}
	for i := 0; i < 10; i++ {
		tests = append(tests, generateCase(rng))
	}
	for i, t := range tests {
		if err := runCase(bin, t); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
