package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
100
3 6
1 0 4
2 -5 4
2 1
3 2
2 4
2 5
2 4
2 1
2 0
5 1
-3 0
2 1
1 9
3
5
2 0
3 5
1 4
2 1
3 6
5 2 1
3 -5 2
2 1
3 2
1 8
2
-3
2 7
2 4
-1 5
2 1
5 10
4 1 5 2 3
1 5 -4 -5 4
2 1
3 2
4 1
5 2
2 10
3 3
5 4
2 1
1 6
5
4
4 10
5 0 1 3
-4 -1 -3 2
2 1
3 1
4 1
4 5
2 3 0 3
-2 3 5 -4
2 1
3 1
4 2
4 5
0 1 0 5
-5 5 3 4
2 1
3 1
4 1
5 10
1 2 2 5 1
-4 2 1 5 -4
2 1
3 2
4 2
5 1
3 2
5 4 5
5 0 -4
2 1
3 2
1 0
0
-2
3 8
2 2 4
-5 4 5
2 1
3 2
4 5
4 1 1 3
4 -1 -5 -3
2 1
3 2
4 2
3 5
5 0 2
4 -5 -5
2 1
3 1
2 9
2 2
1 3
2 1
3 1
3 5 1
-5 -1 -3
2 1
3 2
4 5
2 3 0 0
3 2 2 0
2 1
3 1
4 2
1 7
3
-5
3 5
5 5 1
-3 5 4
2 1
3 1
1 1
1
-2
1 6
0
-4
4 8
4 2 3 3
4 5 -2 1
2 1
3 2
4 1
3 9
1 3 1
0 -4 -4
2 1
3 2
2 1
3 3
-1 -2
2 1
2 9
1 0
-2 2
2 1
3 8
1 0 4
2 -3 4
2 1
3 2
5 7
5 2 3 3 5
5 -2 3 4 -2
2 1
3 2
4 3
5 3
3 0
4 1 2
4 -3 1
2 1
3 2
1 1
4
-5
1 3
1
-5
3 0
3 2 1
-3 5 2
2 1
3 2
5 8
0 4 0 5 4
4 -4 1 -2 -1
2 1
3 2
4 2
5 2
1 10
0
-3
3 8
4 2 2
-4 2 -1
2 1
3 2
4 6
0 1 5 1
-2 -1 0 -5
2 1
3 2
4 2
2 7
4 5
-4 5
2 1
3 6
0 4 3
1 2 -5
2 1
3 2
2 0
0 4
4 -3
2 1
1 8
5
0
2 6
3 0
-5 4
2 1
5 10
2 5 0 5 5
4 -1 -3 1 -1
2 1
3 1
4 1
5 4
4 5
1 3 2 5
-4 -5 -5 2
2 1
3 1
4 3
5 9
1 1 0 5 4
3 1 3 -1 -4
2 1
3 2
4 3
5 4
1 1
3
-4
1 6
1
-5
4 6
5 3 0 3
0 -1 -4 0
2 1
3 1
4 2
1 5
2
-3
1 3
2
-4
5 2
1 0 1 5 5
-4 -5 -1 0 -5
2 1
3 1
4 1
5 4
1 7
2
-1
2 0
1 2
0 2
2 1
3 8
5 2 1
4 -4 -4
2 1
3 1
4 2
1 1 2 4
-2 -2 -3 -1
2 1
3 2
4 3
1 2
4
-5
4 1
5 0 1 3
-1 3 1 -3
2 1
3 2
4 3
3 1
1 3 5
0 5 3
2 1
3 2
4 0
3 5 2 3
-2 0 -1 2
2 1
3 1
4 1
3 1
4 4 5
-3 2 1
2 1
3 2
4 2
1 3 2 4
-3 0 2 5
2 1
3 2
4 1
3 0
5 3 4
2 -5 -2
2 1
3 1
3 8
4 1 3
2 -4 5
2 1
3 1
5 6
2 5 0 0 2
5 -5 -5 -1 1
2 1
3 2
4 1
5 3
3 4
5 1 4
-4 -5 -4
2 1
3 2
5 5
0 4 1 1 0
1 -1 -1 3 -3
2 1
3 1
4 2
5 4
3 4
3 2 4
5 -3 -3
2 1
3 1
4 6
4 3 1 4
5 -1 0 5
2 1
3 2
4 1
4 7
5 4 2 3
5 -5 2 -1
2 1
3 2
4 1
5 3
0 2 3 3 0
3 -4 5 -4 5
2 1
3 1
4 2
5 1
1 9
0
-1
3 3
1 4 2
-2 -4 1
2 1
3 2
4 2
2 3 5 5
1 -3 2 -3
2 1
3 1
4 1
2 7
2 3
1 2
2 1
2 3
3 1
4 -5
2 1
1 3
5
-4
2 5
0 5
5 5
2 1
2 9
2 4
-4 3
2 1
3 6
3 0 5
3 5 5
2 1
3 2
4 4
5 3 1 2
-1 -5 -5 -5
2 1
3 2
4 1
3 10
0 1 0
1 5 -2
2 1
3 1
4 3
2 4 0 4
-4 0 0 3
2 1
3 2
4 2
1 8
0
-2
3 1
1 4 2
-2 -2 -1
2 1
3 2
5 6
2 3 2 5 1
-5 -1 3 -4 -5
2 1
3 2
4 3
5 4
1 6
3
2
4 1
0 0 1 0
-3 1 -2 2
2 1
3 2
4 3
4 0
1 1 3 1
-3 -1 0 0
2 1
3 1
4 3
3 9
4 1 5
-1 2 3
2 1
3 2
3 3
0 0 4
-4 -3 1
2 1
3 1
3 10
0 5 4
3 1 -5
2 1
3 2
3 1
5 4 2
-2 5 3
2 1
3 1
2 1
4 2
5 0
2 1
3 10
3 2 4
-3 -3 -5
2 1
3 2
5 10
0 1 3 1 1
3 -4 -3 -2 2
2 1
3 1
4 3
5 2
2 6
2 4
4 -3
2 1
`

type testCase struct {
	n     int
	m     int64
	p     []int64
	h     []int64
	edges [][2]int
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(tc testCase) bool {
	n := tc.n
	adj := make([][]int, n)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	order := make([]int, 0, n)
	stack := []int{0}
	parent[0] = -1
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}

	sumPop := make([]int64, n)
	good := make([]int64, n)
	ok := true
	for i := len(order) - 1; i >= 0 && ok; i-- {
		u := order[i]
		sum := tc.p[u]
		sumGood := int64(0)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			sum += sumPop[v]
			sumGood += good[v]
		}
		sumPop[u] = sum
		if (sum+tc.h[u])&1 != 0 {
			ok = false
			break
		}
		g := (sum + tc.h[u]) / 2
		if g < 0 || g > sum {
			ok = false
			break
		}
		if sumGood > g {
			ok = false
			break
		}
		if g-sumGood > tc.p[u] {
			ok = false
			break
		}
		good[u] = g
	}
	return ok
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, err
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+2 > len(fields) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, err
		}
		mVal, err := strconv.ParseInt(fields[pos+1], 10, 64)
		if err != nil {
			return nil, err
		}
		pos += 2
		if pos+n+n+(n-1)*2 > len(fields) {
			return nil, fmt.Errorf("case %d missing values", i+1)
		}
		p := make([]int64, n)
		for j := 0; j < n; j++ {
			v, err := strconv.ParseInt(fields[pos+j], 10, 64)
			if err != nil {
				return nil, err
			}
			p[j] = v
		}
		pos += n
		h := make([]int64, n)
		for j := 0; j < n; j++ {
			v, err := strconv.ParseInt(fields[pos+j], 10, 64)
			if err != nil {
				return nil, err
			}
			h[j] = v
		}
		pos += n
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			u, err := strconv.Atoi(fields[pos+2*j])
			if err != nil {
				return nil, err
			}
			v, err := strconv.Atoi(fields[pos+2*j+1])
			if err != nil {
				return nil, err
			}
			edges[j] = [2]int{u - 1, v - 1}
		}
		pos += 2 * (n - 1)
		cases = append(cases, testCase{n: n, m: mVal, p: p, h: h, edges: edges})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("unexpected extra data after %d cases", len(cases))
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range tc.h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		got, err := run(bin, buildInput(tc))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		want := "NO"
		if solve(tc) {
			want = "YES"
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
