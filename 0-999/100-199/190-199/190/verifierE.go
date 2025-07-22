package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func solve(n, m int, edges [][2]int) string {
	neighbors := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		neighbors[u] = append(neighbors[u], v)
		neighbors[v] = append(neighbors[v], u)
	}

	num := make([]int, n+2)
	bel := make([]int, n+2)
	wz := make([]int, n+2)
	dl := make([]int, n+2)
	for i := 1; i <= n; i++ {
		dl[i] = i
		wz[i] = i
	}

	top, now := 0, 0
	for i := 1; i <= n; i++ {
		if now < i {
			now = i
			top++
		}
		x := dl[i]
		bel[x] = top
		num[top]++
		td := n
		for _, j := range neighbors[x] {
			if wz[j] <= now {
				continue
			}
			if wz[j] > td {
				continue
			}
			y := dl[td]
			wz[y], wz[j] = wz[j], wz[y]
			dl[wz[y]], dl[wz[j]] = y, j
			td--
		}
		now = td
	}

	a := make([]int, n)
	for i := 1; i <= n; i++ {
		a[i-1] = i
	}
	sort.Slice(a, func(i, j int) bool {
		if bel[a[i]] != bel[a[j]] {
			return bel[a[i]] < bel[a[j]]
		}
		return a[i] < a[j]
	})

	var out bytes.Buffer
	fmt.Fprintln(&out, top)
	idx := 0
	for comp := 1; comp <= top; comp++ {
		fmt.Fprintf(&out, "%d", num[comp])
		for k := 0; k < num[comp]; k++ {
			fmt.Fprintf(&out, " %d", a[idx])
			idx++
		}
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func generateCases() []testCase {
	rand.Seed(5)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(8) + 1
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges + 1)
		edges := make([][2]int, 0, m)
		used := map[[2]int]bool{}
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			e := [2]int{u, v}
			if used[e] {
				continue
			}
			used[e] = true
			edges = append(edges, e)
		}
		buf := bytes.Buffer{}
		fmt.Fprintf(&buf, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&buf, "%d %d\n", e[0], e[1])
		}
		cases[t] = testCase{input: buf.String(), expected: solve(n, m, edges)}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE.go <binary>")
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
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
