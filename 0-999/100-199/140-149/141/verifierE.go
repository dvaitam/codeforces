package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `1 3 1 1 E 1 1 S 1 1 E
1 2 1 1 S 1 1 S
1 3 1 1 S 1 1 E 1 1 E
5 6 2 2 E 3 1 E 1 3 E 5 2 E 4 5 E 4 4 S
1 3 1 1 S 1 1 E 1 1 S
5 6 1 4 S 4 3 S 3 4 E 5 2 E 1 1 S 3 5 S
1 3 1 1 E 1 1 E 1 1 E
1 3 1 1 S 1 1 E 1 1 E
1 3 1 1 E 1 1 S 1 1 E
3 5 3 1 E 1 2 S 1 1 S 1 3 S 3 3 S
5 5 5 4 S 3 1 S 5 3 E 2 4 S 2 4 E
3 3 1 2 E 1 3 E 1 2 S
1 1 1 1 S
5 6 2 4 E 2 3 S 3 5 S 5 4 S 4 4 S 4 2 S
3 5 3 2 E 2 3 S 3 3 E 2 2 E 1 2 S
1 4 1 1 S 1 1 E 1 1 S 1 1 S
1 5 1 1 E 1 1 S 1 1 E 1 1 S 1 1 S
3 7 2 3 S 1 2 S 3 1 E 2 1 S 3 2 E 1 2 E 1 2 E
1 1 1 1 E
3 6 1 3 S 3 2 E 2 1 S 2 1 S 3 1 E 2 2 E
5 5 2 1 E 1 5 S 4 5 E 5 1 S 4 3 E
1 2 1 1 E 1 1 E
5 9 4 4 S 3 5 E 3 2 E 1 5 S 5 2 E 3 2 E 5 1 S 4 2 E 3 3 S
3 3 3 1 S 2 1 S 1 3 E
5 9 1 1 S 4 2 E 4 5 S 4 2 E 3 2 E 5 5 E 4 5 E 3 3 E 4 5 S
3 4 3 2 E 2 2 S 1 3 S 2 3 E
5 7 1 3 E 2 4 S 2 2 S 3 2 S 2 4 E 1 5 S 2 1 S
1 3 1 1 E 1 1 S 1 1 S
5 5 5 3 S 5 2 E 2 2 E 4 5 S 4 2 S
3 7 2 3 E 2 2 S 2 1 E 1 2 S 2 1 S 1 1 S 1 1 S
3 7 1 2 E 2 1 S 1 2 E 2 3 E 3 3 S 3 1 E 3 3 E
3 5 1 1 E 1 3 S 3 2 E 2 1 S 3 1 E
5 9 5 2 S 2 1 E 5 2 E 1 5 S 1 2 E 5 1 S 1 4 S 2 3 S 3 3 E
5 9 4 1 S 4 4 S 5 5 E 5 3 S 4 3 S 2 1 E 4 3 E 4 4 E 5 2 E
3 7 1 1 E 2 3 E 3 3 S 3 2 S 1 3 S 3 1 E 3 1 E
1 3 1 1 S 1 1 S 1 1 E
3 7 3 3 S 3 1 S 3 3 S 1 3 S 2 1 E 3 2 E 1 1 S
1 1 1 1 E
5 8 1 2 S 4 4 E 5 3 E 4 4 S 1 2 S 3 2 E 2 3 S 5 3 E
3 7 3 1 S 1 2 S 3 1 S 1 1 S 2 3 E 3 2 S 2 3 E
1 2 1 1 E 1 1 E
5 6 4 1 E 2 1 E 3 5 S 1 4 S 5 1 S 5 2 E
5 5 1 4 E 2 5 S 5 1 S 2 5 E 2 4 E
5 7 2 3 E 3 5 E 5 1 S 3 5 S 4 4 S 1 5 S 2 2 S
5 5 5 3 S 1 5 S 5 2 E 5 1 E 4 3 S
5 8 3 5 S 2 3 E 1 1 E 5 5 E 1 2 E 1 4 S 3 2 S 2 2 S
1 1 1 1 E
1 1 1 1 S
5 9 1 1 S 2 4 S 1 5 E 5 3 E 2 5 S 5 5 S 5 1 S 1 3 E 3 5 E
1 5 1 1 S 1 1 E 1 1 E 1 1 S 1 1 E
3 5 3 1 E 3 2 E 2 3 S 2 1 S 2 3 S
1 4 1 1 S 1 1 S 1 1 S 1 1 E
1 4 1 1 S 1 1 S 1 1 E 1 1 S
3 3 3 1 E 2 3 E 3 1 S
3 5 1 2 S 2 2 E 3 3 S 2 1 S 2 2 S
1 1 1 1 S
1 4 1 1 E 1 1 E 1 1 E 1 1 E
3 5 1 3 S 2 2 S 2 3 S 3 3 S 3 1 E
1 5 1 1 E 1 1 E 1 1 E 1 1 S 1 1 E
1 3 1 1 E 1 1 E 1 1 E
5 9 2 5 S 1 4 E 2 1 S 2 3 S 3 3 E 5 2 E 1 1 S 4 1 S 2 4 S
5 9 2 2 S 1 2 S 2 1 E 3 3 E 5 4 S 3 1 E 3 1 S 3 1 E 1 2 S
1 4 1 1 S 1 1 E 1 1 S 1 1 S
1 4 1 1 E 1 1 E 1 1 E 1 1 E
5 5 2 5 S 2 4 S 2 5 S 5 1 S 3 1 S
5 8 1 5 S 3 3 S 2 1 E 5 4 E 5 3 S 1 4 E 5 2 E 5 3 S
3 4 1 1 E 3 1 S 1 3 E 3 2 S
5 6 3 5 E 4 3 E 1 4 E 3 3 S 1 3 E 5 5 E
5 8 1 3 E 5 1 E 5 1 S 4 5 S 4 4 S 1 3 S 2 4 S 2 4 E
3 6 2 2 E 1 3 S 3 3 S 3 1 S 3 2 E 2 1 S
3 3 3 3 E 2 2 S 2 2 E
3 7 1 2 S 1 2 E 2 1 E 3 2 E 3 1 S 1 2 S 1 2 S
5 7 5 2 S 4 4 E 4 2 S 1 2 E 3 5 E 1 3 E 5 4 S
1 3 1 1 E 1 1 S 1 1 E
5 5 2 5 S 1 4 E 4 1 E 4 4 E 5 2 E
5 6 3 3 E 1 4 S 3 5 E 3 5 E 5 3 S 3 5 E
3 7 2 1 S 1 3 E 3 1 E 2 2 E 3 2 S 3 3 E 1 2 E
3 4 2 2 S 1 2 S 3 2 S 1 2 S
1 3 1 1 E 1 1 S 1 1 S
5 7 3 1 E 5 5 S 5 4 E 5 5 E 2 2 E 4 2 E 4 5 S
1 2 1 1 E 1 1 E
1 5 1 1 S 1 1 E 1 1 E 1 1 E 1 1 S
5 6 5 1 E 1 1 S 2 5 S 2 5 E 2 1 E 3 3 E
3 4 3 1 E 2 3 E 3 3 E 1 3 S
5 8 3 1 E 1 3 S 5 5 E 1 4 E 4 5 E 5 4 E 5 1 S 5 2 E
5 7 5 3 E 1 3 S 4 2 E 3 1 S 4 3 S 5 5 E 1 2 S
5 5 4 2 E 5 1 S 2 5 S 2 2 S 3 2 E
1 3 1 1 E 1 1 E 1 1 S
3 6 2 3 E 3 2 E 3 2 S 2 1 E 1 3 S 2 2 E
1 3 1 1 E 1 1 S 1 1 S
5 7 4 4 S 5 3 S 2 1 S 4 2 S 2 1 E 3 2 S 2 5 S
3 7 2 2 S 2 2 S 2 2 S 2 1 E 1 1 S 1 2 E 2 3 E
1 2 1 1 E 1 1 E
5 8 3 5 S 4 3 E 4 3 S 4 4 S 2 1 E 5 3 E 2 5 E 2 3 E
3 7 1 1 E 1 1 E 1 3 S 3 2 S 1 2 E 3 2 S 1 3 S
1 1 1 1 S
5 6 2 3 E 1 2 E 2 5 E 1 5 E 5 3 E 2 4 S
5 8 4 5 S 2 3 E 5 4 S 1 1 S 3 2 S 4 5 S 4 4 E 5 3 S
3 4 2 3 E 1 3 S 1 1 E 1 2 E
3 3 1 2 E 3 1 S 2 2 E`

type DSU struct{ p []int }

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(x, y int) {
	fx := d.Find(x)
	fy := d.Find(y)
	if fx != fy {
		d.p[fx] = fy
	}
}

type Edge struct {
	u, v int
	s    bool
}

type testCase struct {
	n     int
	m     int
	edges []Edge
}

func solveCase(tc testCase) string {
	if tc.n%2 == 0 {
		return "-1"
	}
	initialUsd := make([]bool, tc.m)
	usd := make([]bool, tc.m)
	dsu := NewDSU(tc.n)
	num := 0
	for i := 0; i < tc.m; i++ {
		e := tc.edges[i]
		if dsu.Find(e.u) != dsu.Find(e.v) {
			dsu.Union(e.u, e.v)
			initialUsd[i] = true
			if e.s {
				num++
			}
		}
	}
	copy(usd, initialUsd)
	if num*2+1 != tc.n {
		dsu = NewDSU(tc.n)
		if num*2+1 < tc.n {
			newUsd := make([]bool, tc.m)
			newNum := 0
			for i := 0; i < tc.m; i++ {
				if tc.edges[i].s && initialUsd[i] {
					dsu.Union(tc.edges[i].u, tc.edges[i].v)
					newUsd[i] = true
					newNum++
				}
			}
			for i := 0; i < tc.m && newNum*2+1 != tc.n; i++ {
				if tc.edges[i].s && dsu.Find(tc.edges[i].u) != dsu.Find(tc.edges[i].v) {
					dsu.Union(tc.edges[i].u, tc.edges[i].v)
					newUsd[i] = true
					newNum++
				}
			}
			for i := 0; i < tc.m; i++ {
				if !tc.edges[i].s && dsu.Find(tc.edges[i].u) != dsu.Find(tc.edges[i].v) {
					dsu.Union(tc.edges[i].u, tc.edges[i].v)
					newUsd[i] = true
				}
			}
			usd = newUsd
			num = newNum
		} else {
			newUsd := make([]bool, tc.m)
			newNum := num
			for i := 0; i < tc.m; i++ {
				if !tc.edges[i].s && initialUsd[i] {
					dsu.Union(tc.edges[i].u, tc.edges[i].v)
					newUsd[i] = true
				}
			}
			for i := 0; i < tc.m && newNum*2+1 != tc.n; i++ {
				if !tc.edges[i].s && dsu.Find(tc.edges[i].u) != dsu.Find(tc.edges[i].v) {
					dsu.Union(tc.edges[i].u, tc.edges[i].v)
					newUsd[i] = true
					newNum--
				}
			}
			for i := 0; i < tc.m; i++ {
				if tc.edges[i].s && dsu.Find(tc.edges[i].u) != dsu.Find(tc.edges[i].v) {
					dsu.Union(tc.edges[i].u, tc.edges[i].v)
					newUsd[i] = true
				}
			}
			usd = newUsd
			num = newNum
		}
	}
	if num*2+1 == tc.n {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n-1))
		cnt := 0
		for i := 0; i < tc.m && cnt < tc.n-1; i++ {
			if usd[i] {
				sb.WriteString(fmt.Sprintf("%d ", i+1))
				cnt++
			}
		}
		return strings.TrimSpace(sb.String())
	}
	return "-1"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err1 := strconv.Atoi(parts[0])
		m, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: parse n/m: %v %v", idx+1, err1, err2)
		}
		if len(parts) != 2+3*m {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", idx+1, 2+3*m, len(parts))
		}
		edges := make([]Edge, m)
		pos := 2
		for i := 0; i < m; i++ {
			u, errU := strconv.Atoi(parts[pos])
			v, errV := strconv.Atoi(parts[pos+1])
			typ := parts[pos+2]
			if errU != nil || errV != nil {
				return nil, fmt.Errorf("line %d: parse edge %d: %v %v", idx+1, i+1, errU, errV)
			}
			edges[i] = Edge{u: u, v: v, s: typ == "S"}
			pos += 3
		}
		cases = append(cases, testCase{n: n, m: m, edges: edges})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			t := 'E'
			if e.s {
				t = 'S'
			}
			fmt.Fprintf(&input, "%d %d %c\n", e.u, e.v, t)
		}

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solveCase(tc)
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
