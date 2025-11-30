package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `
100
8 10
3 8
1 2
6 7
1 5
1 4
1 7
1 8
7 8
1 3
1 6
6 4 3 5
7 7
1 6
1 2
1 5
1 7
1 3
1 4
4 6
2 5 2 2
8 8
1 7
1 5
1 2
3 5
1 3
1 6
1 4
1 8
3 1 5 6
4 6
2 4
1 4
1 2
3 4
1 3
2 3
4 3 3 2
6 6
1 3
4 5
1 2
1 4
1 6
1 5
1 3 1 2
4 4
1 3
1 4
2 3
1 2
2 1 3 3
6 5
1 3
1 5
1 4
1 6
1 2
2 6 1 4
3 2
1 3
1 2
3 2 1 2
2 1
1 2
1 2 1 1
3 2
1 2
1 3
1 3 2 1
4 5
2 3
1 2
1 4
1 3
2 4
1 3 1 1
4 5
3 4
1 4
2 4
1 2
1 3
3 4 2 2
8 9
1 3
1 5
1 4
1 6
1 2
1 7
4 7
3 7
1 8
6 2 2 2
7 9
2 7
1 3
1 7
1 2
2 5
2 6
1 5
1 4
1 6
3 7 4 1
8 10
1 7
4 8
6 7
1 6
1 2
1 5
1 8
5 6
1 3
1 4
2 7 3 2
7 8
1 4
1 2
1 3
1 6
1 5
3 6
1 7
3 7
2 4 5 3
2 1
1 2
1 2 1 1
8 9
1 7
1 2
4 7
1 6
1 8
1 5
1 4
1 3
4 5
2 5 4 4
4 5
1 4
2 4
1 3
3 4
1 2
3 2 3 3
3 2
1 2
1 3
3 1 1 1
3 2
1 3
1 2
2 1 1 2
6 7
3 4
1 6
3 5
1 2
1 4
1 3
1 5
5 4 2 1
5 7
1 4
3 4
4 5
3 5
1 5
1 3
1 2
1 4 2 4
7 9
2 5
5 7
1 5
4 7
1 4
1 3
1 6
1 2
1 7
6 7 4 1
8 7
1 5
1 3
1 4
1 2
1 8
1 6
1 7
4 3 5 5
5 7
1 5
1 2
2 3
4 5
1 4
1 3
2 5
4 3 3 4
5 7
1 5
4 5
1 3
2 5
1 2
1 4
3 4
2 4 1 2
8 9
5 7
1 8
1 7
1 4
1 5
1 6
1 2
1 3
2 6
5 7 7 2
5 6
1 2
2 5
3 4
1 3
1 4
1 5
1 5 1 4
4 3
1 3
1 4
1 2
4 3 2 1
3 3
1 2
1 3
2 3
1 3 2 1
2 1
1 2
1 2 1 1
7 9
4 6
1 7
1 6
1 3
5 7
1 4
2 6
1 5
1 2
2 6 6 6
2 1
1 2
2 1 1 1
8 8
1 7
1 2
1 5
4 6
1 4
1 8
1 6
1 3
3 7 2 3
6 6
1 4
1 3
3 5
1 6
1 2
1 5
3 5 4 2
6 8
1 6
1 3
3 6
5 6
1 4
4 6
1 5
1 2
3 4 1 2
8 7
1 8
1 7
1 6
1 2
1 4
1 3
1 5
8 1 2 3
2 1
1 2
1 2 1 1
7 9
6 7
1 2
1 5
1 6
3 6
1 4
1 7
1 3
5 7
5 1 1 1
8 9
1 8
1 7
1 2
1 3
6 8
1 4
1 5
5 6
1 6
2 4 4 5
5 5
1 2
1 4
2 3
1 3
1 5
5 3 4 1
5 6
4 5
3 5
1 3
1 5
1 4
1 2
5 1 1 4
2 1
1 2
1 2 1 1
8 9
1 2
1 4
1 6
1 5
5 8
1 8
1 3
1 7
2 3
7 8 4 4
7 7
1 2
1 5
1 4
1 6
1 3
2 5
1 7
5 1 6 5
8 9
1 7
3 8
1 4
1 8
1 5
2 5
1 3
1 2
1 6
4 8 2 3
6 5
1 4
1 3
1 5
1 2
1 6
1 2 5 3
3 2
1 2
1 3
3 2 2 2
4 6
2 3
1 3
3 4
2 4
1 4
1 2
4 3 2 1
8 10
4 5
1 2
1 8
1 6
1 7
1 4
1 3
2 3
3 6
1 5
5 8 5 7
3 3
1 2
2 3
1 3
3 2 2 1
2 1
1 2
1 2 1 1
3 3
1 3
1 2
2 3
3 2 1 1
5 6
1 5
3 4
1 4
1 2
2 4
1 3
3 2 3 4
4 4
2 4
1 2
1 3
1 4
2 4 1 1
6 5
1 4
1 2
1 6
1 5
1 3
6 4 1 5
8 7
1 5
1 2
1 8
1 6
1 4
1 3
1 7
4 1 3 6
7 7
1 5
1 6
1 2
1 4
1 3
3 7
1 7
6 1 1 2
2 1
1 2
1 2 1 1
4 6
1 2
1 4
3 4
1 3
2 3
2 4
2 3 2 3
8 10
3 5
1 8
2 6
1 7
1 6
1 3
1 2
5 7
1 4
1 5
7 8 5 1
6 8
1 6
4 5
1 5
1 3
2 4
1 4
1 2
3 6
5 3 2 4
4 5
3 4
1 3
2 3
1 4
1 2
4 1 3 3
6 7
1 6
1 3
1 5
3 4
4 6
1 4
1 2
6 2 5 5
4 6
3 4
1 3
2 4
1 4
2 3
1 2
2 1 2 3
8 7
1 3
1 7
1 6
1 8
1 5
1 2
1 4
2 5 4 3
8 8
1 6
1 5
1 4
1 7
1 2
1 3
7 8
1 8
6 4 1 5
5 5
1 3
3 4
1 2
1 5
1 4
5 2 3 2
8 10
4 5
1 7
1 6
1 5
1 3
1 4
1 8
2 3
3 7
1 2
7 6 7 1
6 8
1 3
1 5
5 6
2 4
1 4
1 6
3 6
1 2
2 1 5 4
2 1
1 2
1 2 1 1
8 9
1 2
1 5
1 6
1 7
4 5
5 8
1 8
1 4
1 3
5 8 4 2
8 10
1 6
3 6
1 4
1 5
3 4
1 8
1 2
1 7
1 3
2 3
8 3 7 4
6 7
4 6
1 6
1 3
1 4
2 6
1 5
1 2
5 2 3 1
5 4
1 4
1 2
1 3
1 5
1 4 3 2
5 4
1 4
1 3
1 5
1 2
5 1 3 4
8 10
1 5
2 7
1 2
5 7
1 3
3 4
1 6
1 4
1 7
1 8
4 5 3 1
5 7
4 5
1 5
1 2
2 5
3 4
1 3
1 4
4 5 3 2
5 4
1 5
1 4
1 2
1 3
1 3 1 1
2 1
1 2
1 2 1 1
5 7
3 4
1 4
1 2
1 5
2 3
1 3
4 5
4 3 1 4
5 5
1 3
1 2
1 4
2 5
1 5
3 5 3 1
6 8
1 6
1 3
3 5
3 6
1 2
1 5
2 4
1 4
4 1 4 5
5 5
1 3
1 2
1 5
1 4
3 5
3 5 2 3
8 10
6 7
1 5
1 8
1 3
6 8
1 6
1 2
3 5
1 4
1 7
7 4 1 6
7 6
1 2
1 3
1 7
1 6
1 4
1 5
6 1 4 3
3 3
1 2
1 3
2 3
3 1 2 2
3 3
1 2
2 3
1 3
2 3 2 2
8 7
1 7
1 3
1 6
1 8
1 5
1 2
1 4
8 7 7 2
3 2
1 3
1 2
1 3 2 1
5 7
1 4
1 3
3 5
3 4
1 2
4 5
1 5
3 2 4 1
5 4
1 4
1 2
1 3
1 5
2 3 2 3
8 8
1 3
1 7
1 4
1 6
7 8
1 8
1 2
1 5
8 5 3 7
8 9
1 2
1 6
1 7
1 5
4 5
5 7
1 3
1 4
1 8
5 4 5 2
8 7
1 5
1 2
1 6
1 4
1 8
1 7
1 3
2 6 4 4
5 7
3 4
1 2
2 5
1 3
1 5
1 4
2 3
1 4 3 4
3 3
1 3
2 3
1 2
2 1 2 1
7 8
1 4
1 6
5 7
1 3
5 6
1 7
1 2
1 5
1 5 2 2
3 3
1 2
2 3
1 3
1 2 1 2

`

type Pair struct{ u, v int }

type testCase struct {
	n, m   int
	edges  []Pair
	s, t   int
	ds, dt int
}

func parseCases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcases))
	scan.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scan.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("read test count: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read n: %w", i+1, err)
		}
		m, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read m: %w", i+1, err)
		}
		edges := make([]Pair, m)
		for j := 0; j < m; j++ {
			u, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d: read u: %w", i+1, j+1, err)
			}
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d edge %d: read v: %w", i+1, j+1, err)
			}
			edges[j] = Pair{u, v}
		}
		s, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read s: %w", i+1, err)
		}
		tval, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read t: %w", i+1, err)
		}
		ds, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read ds: %w", i+1, err)
		}
		dt, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: read dt: %w", i+1, err)
		}
		cases = append(cases, testCase{n: n, m: m, edges: edges, s: s, t: tval, ds: ds, dt: dt})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func referenceSolve(tc testCase) string {
	n, s, t := tc.n, tc.s, tc.t
	ds, dt := tc.ds, tc.dt
	edges := tc.edges
	if s > t {
		s, t = t, s
		ds, dt = dt, ds
	}

	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(u int) int {
		if parent[u] != u {
			parent[u] = find(parent[u])
		}
		return parent[u]
	}

	ans := make([]Pair, 0, n)
	for _, e := range edges {
		u, v := e.u, e.v
		if u != s && v != s && u != t && v != t {
			ru, rv := find(u), find(v)
			if ru != rv {
				parent[ru] = rv
				ans = append(ans, Pair{u, v})
			}
		}
	}

	a := make([]int, n+1)
	b := make([]int, n+1)
	flag := false
	for _, e := range edges {
		u, v := e.u, e.v
		if u == s && v == t {
			flag = true
		} else {
			if u == s && v != t {
				a[find(v)] = v
			}
			if v == s && u != t {
				a[find(u)] = u
			}
			if u == t && v != s {
				b[find(v)] = v
			}
			if v == t && u != s {
				b[find(u)] = u
			}
		}
	}

	x := make([]Pair, 0)
	y := make([]Pair, 0)
	z := make([]Pair, 0)
	for i := 1; i <= n; i++ {
		ri := find(i)
		if ri == i && i != s && i != t {
			av := a[i]
			bv := b[i]
			if av != 0 && bv != 0 {
				x = append(x, Pair{av, bv})
			} else if av != 0 {
				y = append(y, Pair{s, av})
			} else if bv != 0 {
				z = append(z, Pair{bv, t})
			} else {
				return "No"
			}
		}
	}

	ds -= len(y)
	dt -= len(z)
	ans = append(ans, y...)
	ans = append(ans, z...)

	if len(x) > 0 {
		ds--
		dt--
		last := x[len(x)-1]
		ans = append(ans, Pair{last.u, s})
		ans = append(ans, Pair{last.v, t})
		x = x[:len(x)-1]
		take := ds
		if take < 0 {
			take = 0
		}
		if take > len(x) {
			take = len(x)
		}
		ds -= take
		for i := 0; i < take; i++ {
			ans = append(ans, Pair{x[i].u, s})
		}
		rem := len(x) - take
		dt -= rem
		for i := take; i < len(x); i++ {
			ans = append(ans, Pair{x[i].v, t})
		}
	} else if flag {
		ans = append(ans, Pair{s, t})
		ds--
		dt--
	} else {
		return "No"
	}
	if ds < 0 || dt < 0 {
		return "No"
	}

	var sb strings.Builder
	sb.WriteString("Yes\n")
	for _, p := range ans {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.u, p.v))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e.u, e.v)
		}
		fmt.Fprintf(&input, "%d %d %d %d\n", tc.s, tc.t, tc.ds, tc.dt)

		expected := referenceSolve(tc)
		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
