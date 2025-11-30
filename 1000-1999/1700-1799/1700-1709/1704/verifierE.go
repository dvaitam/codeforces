package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 1 0 3 2 1
4 4 0 3 0 0 3 1 3 2 3 4 1 4
2 2 1 1 1 2 2 1
2 0 3 2 
3 3 3 0 1 3 1 1 3 2 1
3 2 1 1 2 1 3 2 1
3 5 0 3 2 2 1 3 1 2 3 3 2 1 3
3 3 0 1 1 2 3 1 2 2 1
4 5 0 1 1 3 1 2 4 1 4 3 4 2 1 3
3 4 3 0 1 2 3 3 2 2 1 3 1
4 11 2 0 0 2 2 4 1 2 2 1 3 4 4 3 3 1 1 4 4 2 2 3 3 2 1 3
4 12 3 2 0 0 1 3 2 4 1 2 2 1 3 4 4 3 3 1 1 4 4 2 2 3 3 2 4 1
4 7 1 3 3 0 1 3 1 2 4 3 3 1 2 3 3 2 4 1
2 1 2 1 1 2
3 1 1 2 0 1 2
4 3 3 1 1 1 3 1 1 2 1 4
3 4 1 0 3 3 2 1 2 1 3 2 1
3 3 3 2 3 2 3 3 2 1 2
2 2 1 2 1 2 2 1
2 0 0 3 
2 0 0 0 
3 4 3 0 0 3 2 1 2 1 3 2 1
3 4 2 3 3 3 1 3 2 1 2 1 3
2 2 0 1 1 2 2 1
2 2 1 3 1 2 2 1
2 2 3 3 1 2 2 1
4 0 2 2 0 3 
3 4 0 2 3 3 1 3 2 1 3 2 1
3 5 0 1 2 1 2 3 1 2 3 3 2 1 3
4 1 3 2 3 3 1 3
2 2 2 3 1 2 2 1
3 1 2 0 0 2 3
4 3 1 3 1 2 3 1 4 1 4 2
2 2 0 2 1 2 2 1
4 7 0 3 1 2 2 4 2 1 4 3 4 1 1 4 3 2 1 3
4 12 2 0 1 0 2 4 1 2 3 4 2 1 4 3 3 1 4 1 4 2 1 4 2 3 3 2 1 3
3 2 3 0 1 2 3 1 3
4 10 0 2 3 1 1 3 2 4 2 1 4 3 3 1 4 2 1 4 2 3 3 2 4 1
2 2 1 0 1 2 2 1
4 0 3 1 0 3 
4 6 1 1 1 0 1 2 2 1 3 4 4 3 2 3 3 2
3 1 2 1 1 1 3
2 1 2 2 1 2
3 0 0 1 3 
3 2 2 3 2 3 1 3 2
2 1 0 0 2 1
3 1 3 2 0 1 3
4 9 3 3 3 3 2 4 1 2 3 4 2 1 4 3 3 1 4 2 2 3 4 1
2 1 2 1 1 2
4 10 1 1 0 1 2 4 1 2 2 1 3 4 4 3 3 1 4 2 1 4 3 2 4 1
4 9 3 3 1 0 2 4 1 2 3 4 2 1 4 3 3 1 1 4 4 2 2 3
4 11 2 3 1 1 1 3 2 4 1 2 2 1 3 4 4 3 3 1 1 4 2 3 3 2 4 1
4 7 2 0 2 2 2 4 1 2 3 4 4 3 1 4 4 2 1 3
4 4 0 0 1 3 2 3 3 2 1 2 4 2
2 0 0 1 
2 2 2 2 1 2 2 1
4 6 3 2 1 3 2 1 3 4 4 3 4 1 2 3 1 3
2 0 3 0 
4 3 1 3 3 1 3 1 3 2 1 3
2 2 1 1 1 2 2 1
2 2 0 2 1 2 2 1
2 0 2 2 
3 6 2 2 3 1 2 2 1 3 1 2 3 3 2 1 3
4 11 1 3 3 0 2 4 1 2 2 1 3 4 4 3 3 1 4 1 4 2 1 4 3 2 1 3
2 0 0 1 
3 0 1 2 1 
3 2 3 3 1 3 2 2 1
4 0 0 2 3 1 
4 4 3 1 3 3 2 4 1 3 3 4 4 2
3 6 2 0 0 1 2 2 1 3 1 2 3 3 2 1 3
2 2 3 3 1 2 2 1
3 5 2 1 2 1 2 2 1 3 1 2 3 3 2
4 6 1 2 3 2 1 2 2 1 3 4 4 3 3 1 1 3
3 6 3 0 2 1 2 2 1 3 1 2 3 3 2 1 3
2 2 3 0 1 2 2 1
3 4 2 0 0 2 3 3 2 1 2 2 1
2 0 2 0 
2 2 2 0 1 2 2 1
4 8 3 3 3 3 2 4 3 4 4 3 4 2 1 4 2 3 3 2 1 3
2 0 3 2 
4 4 2 1 1 3 1 2 1 3 1 4 4 3
3 6 1 1 0 1 2 2 1 3 1 2 3 3 2 1 3
2 0 1 1 
4 9 3 0 0 1 2 4 3 4 4 1 3 1 4 3 1 4 4 2 3 2 1 3
4 11 1 1 3 1 2 4 1 2 2 1 3 4 4 3 3 1 1 4 4 2 2 3 3 2 1 3
3 3 0 1 2 3 2 1 2 1 3
2 0 3 2 
3 5 2 2 2 1 2 2 1 3 1 2 3 1 3
2 2 1 3 1 2 2 1
2 1 3 0 2 1
4 12 2 3 3 0 2 4 1 2 2 1 3 4 4 1 3 1 4 3 4 2 1 4 2 3 3 2 1 3
3 1 3 1 3 2 1
4 8 2 1 3 3 2 4 2 1 4 1 4 3 4 2 1 4 3 2 1 3
2 0 0 2 
3 6 2 1 2 1 2 2 1 3 1 2 3 3 2 1 3
3 5 0 1 3 1 2 2 1 2 3 3 2 1 3
3 3 3 1 3 3 1 3 2 2 1
3 5 0 0 1 2 1 3 1 2 3 3 2 1 3
2 0 1 3 
4 6 2 1 2 1 2 4 2 1 4 3 3 1 1 4 3 2`

type pair struct {
	a, b int
}

type testCase struct {
	n, m  int
	a     []int64
	edges []pair
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: invalid", i+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", i+1, err)
		}
		if len(parts) != 2+n+2*m {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, 2+n+2*m, len(parts))
		}
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			v, err := strconv.ParseInt(parts[2+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", i+1, j, err)
			}
			a[j] = v
		}
		edges := make([]pair, m)
		offset := 2 + n
		for j := 0; j < m; j++ {
			x, err := strconv.Atoi(parts[offset+2*j])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge %d u: %v", i+1, j, err)
			}
			y, err := strconv.Atoi(parts[offset+2*j+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse edge %d v: %v", i+1, j, err)
			}
			edges[j] = pair{a: x, b: y}
		}
		cases = append(cases, testCase{n: n, m: m, a: a, edges: edges})
	}
	return cases, nil
}

const mod int64 = 998244353

func solve(tc testCase) int64 {
	n := tc.n
	a := make([]int64, n)
	copy(a, tc.a)
	adj := make([][]int, n)
	outDeg := make([]int, n)
	for _, e := range tc.edges {
		u := e.a - 1
		v := e.b - 1
		adj[u] = append(adj[u], v)
		outDeg[u]++
	}

	ans := 0
	for step := 0; step < n; step++ {
		allZero := true
		inc := make([]int64, n)
		for i := 0; i < n; i++ {
			if a[i] > 0 {
				allZero = false
				a[i]--
				for _, v := range adj[i] {
					inc[v]++
				}
			}
		}
		for i := 0; i < n; i++ {
			a[i] += inc[i]
		}
		if allZero {
			ans = step
			break
		}
		if step == n-1 {
			ans = n
		}
	}
	if ans < n {
		return int64(ans)
	}

	inDeg := make([]int, n)
	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			inDeg[v]++
		}
	}
	queue := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if inDeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	order := make([]int, 0, n)
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		order = append(order, u)
		for _, v := range adj[u] {
			inDeg[v]--
			if inDeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	for _, u := range order {
		for _, v := range adj[u] {
			a[v] = (a[v] + a[u]) % mod
		}
	}
	sink := 0
	for i := 0; i < n; i++ {
		if outDeg[i] == 0 {
			sink = i
			break
		}
	}
	return (int64(ans) + a[sink]) % mod
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

func runBinary(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.a, e.b))
		}

		expect := fmt.Sprint(solve(tc))
		got, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
