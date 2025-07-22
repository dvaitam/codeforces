package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func solve(n int, edges [][2]int) string {
	f := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		f[u] = append(f[u], v)
		f[v] = append(f[v], u)
	}
	a := make([]int, n+1)
	use := make([]bool, n+1)
	a[1] = 1
	var getRem func() bool
	getRem = func() bool {
		for i := 3; i <= n-1; i++ {
			found := false
			for _, u := range f[a[i-1]] {
				if use[u] {
					continue
				}
				ok := false
				for _, v := range f[a[i-2]] {
					if v == u {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				use[u] = true
				a[i] = u
				found = true
				break
			}
			if !found {
				return false
			}
		}
		return true
	}
	check := func() bool {
		t1, t2, t3, t4 := false, false, false, false
		for _, v := range f[a[n]] {
			if v == a[1] {
				t1 = true
			}
			if v == a[2] {
				t2 = true
			}
			if v == a[n-1] {
				t3 = true
			}
			if v == a[n-2] {
				t4 = true
			}
		}
		return t1 && t2 && t3 && t4
	}
	for p := 0; p < len(f[1]); p++ {
		for q := 0; q < len(f[1]); q++ {
			if p == q {
				continue
			}
			u := f[1][p]
			v := f[1][q]
			ok := false
			for _, x := range f[u] {
				if x == v {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
			for i := range use {
				use[i] = false
			}
			use[1], use[u], use[v] = true, true, true
			a[2] = v
			a[n] = u
			if getRem() && check() {
				var out bytes.Buffer
				for i := 1; i <= n; i++ {
					if i > 1 {
						out.WriteByte(' ')
					}
					fmt.Fprintf(&out, "%d", a[i])
				}
				return out.String()
			}
		}
	}
	return "-1"
}

func generateCases() []testCase {
	rand.Seed(3)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(6) + 5
		perm := make([]int, n)
		perm[0] = 1
		p := rand.Perm(n - 1)
		for i, v := range p {
			perm[i+1] = v + 2
		}
		edgeMap := map[[2]int]bool{}
		for i := 0; i < n; i++ {
			a := perm[i]
			b := perm[(i+1)%n]
			if a > b {
				a, b = b, a
			}
			edgeMap[[2]int{a, b}] = true
			a2 := perm[i]
			b2 := perm[(i+2)%n]
			if a2 > b2 {
				a2, b2 = b2, a2
			}
			edgeMap[[2]int{a2, b2}] = true
		}
		edges := make([][2]int, 0, len(edgeMap))
		for e := range edgeMap {
			edges = append(edges, e)
		}
		rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
		var buf bytes.Buffer
		fmt.Fprintln(&buf, n)
		for _, e := range edges {
			fmt.Fprintf(&buf, "%d %d\n", e[0], e[1])
		}
		expected := solve(n, edges)
		cases[t] = testCase{input: buf.String(), expected: expected}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
