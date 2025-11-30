package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesEData = `
4 0 1 0 3 4 2 1 1 1 4 3 1 2
9 0 1 1 1 0 2 1 0 4 3 4 3 5 5 6 2
10 0 1 2 3 4 1 1 3 7 4 1 9 5
1 0 3 1 1 1 1 1 1
8 0 0 0 2 2 0 0 0 4 5 8 6 3 4 2 7 4
8 0 1 0 2 3 5 4 5 5 4 6 2 1 4 5 4 2 6 3
5 0 1 0 0 2 1 3 3
1 0 3 1 1 1 1 1 1
10 0 0 1 2 1 2 3 2 5 9 1 6 1
8 0 0 1 2 2 4 0 7 2 7 4 2 1
1 0 1 1 1
10 0 0 2 3 4 1 2 0 1 8 3 7 4 8 4 4 8
7 0 1 0 1 3 3 1 4 7 2 4 2 1 1 3 3
4 0 0 0 3 3 2 3 1 3 1 4
1 0 4 1 1 1 1 1 1 1 1
6 0 1 2 1 2 2 4 4 1 3 6 6 2 1 4
10 0 0 1 0 1 5 5 7 7 6 4 4 1 4 3 1 10 5 2
7 0 1 0 0 1 1 5 5 3 7 5 7 7 4 5 4 1 1
1 0 5 1 1 1 1 1 1 1 1 1 1
2 0 0 5 2 2 1 2 2 1 2 2 1 1
2 0 1 5 2 2 2 1 1 2 2 2 1 1
6 0 0 0 1 2 3 3 3 5 1 6 2 1
5 0 0 2 0 3 5 4 5 1 5 2 4 3 3 2 2
1 0 1 1 1
10 0 1 2 3 4 4 3 6 2 4 5 6 6 3 7 2 10 3 10 3 5
6 0 0 2 2 4 0 1 4 6
3 0 1 2 3 2 1 2 1 3 1
9 0 0 1 0 1 1 4 7 1 1 3 8
4 0 1 2 3 5 2 4 2 3 3 2 3 4 4 4
6 0 1 1 3 4 0 3 2 5 4 6 5 5
6 0 1 1 0 1 4 1 4 6
6 0 1 0 2 3 1 4 5 1 2 2 1 3 2 1
3 0 1 2 5 3 3 3 2 1 3 1 1 1 1
1 0 3 1 1 1 1 1 1
3 0 1 0 1 3 2
1 0 5 1 1 1 1 1 1 1 1 1 1
5 0 1 2 2 3 3 2 4 1 4 1 3
1 0 3 1 1 1 1 1 1
4 0 0 1 1 4 4 3 2 2 1 4 3 3
7 0 0 2 1 2 4 6 5 3 7 4 6 7 3 7 7 1 2
6 0 0 2 0 4 2 5 4 3 1 2 6 1 4 5 5 5
4 0 0 0 0 3 2 3 1 1 1 2
10 0 1 2 0 0 0 3 2 3 5 2 6 5 7 10
7 0 0 0 3 3 0 6 5 7 5 3 5 6 3 7 2 6 6
7 0 1 0 1 0 3 3 3 3 7 4 6 7 4
6 0 0 1 2 4 0 1 4 3
10 0 1 2 1 4 2 0 2 1 3 3 4 8 2 6 1 8
10 0 0 0 0 1 4 6 1 1 5 3 10 10 3 2 3 1
3 0 0 1 1 3 2
6 0 1 0 1 1 0 1 2 6
1 0 3 1 1 1 1 1 1
10 0 1 2 2 2 3 5 6 0 0 2 3 6 4 5
4 0 1 2 1 5 3 4 3 1 2 2 3 2 2 1
1 0 1 1 1
6 0 1 0 2 3 1 2 4 5 4 3
6 0 0 1 1 3 0 4 2 1 3 3 5 5 1 1
8 0 0 2 0 1 4 1 5 5 1 1 7 7 4 3 1 3 3 5
4 0 1 2 2 5 3 2 3 3 3 3 1 2 3 2
8 0 1 0 0 4 1 1 2 5 3 1 6 2 1 3 3 5 1 5
8 0 1 0 3 2 4 5 0 2 5 7 1 1
6 0 1 0 1 3 5 1 4 1
5 0 0 2 1 1 2 1 1 1 1
5 0 1 0 0 4 3 2 2 5 5 1 1
2 0 0 4 1 1 2 2 2 1 1 1
10 0 0 0 0 2 4 5 6 4 9 3 3 9 6 9 8 4
10 0 1 0 3 4 4 0 3 4 9 5 10 8 6 3 5 3 6 5 6 4
9 0 1 2 2 2 2 6 2 5 1 4 6
4 0 0 1 0 3 1 1 1 1 2 2
4 0 1 2 3 4 2 4 4 3 1 3 1 1
3 0 0 1 1 2 1
7 0 0 2 1 2 2 4 3 6 2 1 3 3 1
7 0 0 1 2 4 3 4 1 4 2
2 0 1 4 1 1 1 2 1 1 1 2
8 0 1 1 2 3 5 2 6 1 6 7
9 0 0 2 3 3 5 2 3 2 4 5 2 6 9 1 9 4 3
8 0 0 2 2 3 3 3 1 4 1 7 6 8 5 4 3 6
6 0 0 1 3 4 1 4 4 6 4 3 5 5 1 6
7 0 1 2 1 4 1 6 1 4 6
10 0 1 0 0 0 1 6 5 2 5 3 8 10 4 5 1 2
1 0 4 1 1 1 1 1 1 1 1
1 0 2 1 1 1 1
6 0 1 1 3 4 1 3 6 6 6 1 5 3
5 0 0 0 2 0 4 1 2 2 1 4 4 3 1
2 0 0 4 2 2 1 1 1 1 1 2
9 0 1 1 2 4 3 1 4 4 5 1 5 7 6 5 2 3 9 4 3
7 0 0 2 1 4 1 4 1 5 6
2 0 1 1 1 2
2 0 0 3 2 1 1 1 2 2
5 0 1 0 0 3 3 5 2 4 5 3 2
8 0 0 0 1 2 5 0 3 2 8 6 3 3
6 0 1 2 3 2 2 1 4 4
5 0 1 0 0 2 4 5 5 4 4 1 3 3 5
2 0 1 5 1 1 2 1 2 2 1 1 2 1
4 0 1 2 1 2 4 4 3 4
5 0 1 1 1 0 5 1 5 2 1 5 4 3 3 1 2
6 0 0 2 3 3 2 2 3 4 3 4
10 0 0 1 1 1 3 3 4 2 7 5 6 6 10 2 3 3 5 10 9 3
2 0 1 5 1 2 1 1 2 2 2 2 1 1
2 0 0 4 2 2 1 2 2 2 1 1
6 0 1 1 2 3 2 1 5 5
`

type testCase struct {
	n       int
	parents []int // 1-based
	queries [][2]int
}

// parseTestcases reads the embedded testcase data.
func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesEData, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		pos := 0
		parseInt := func() (int, error) {
			if pos >= len(fields) {
				return 0, fmt.Errorf("line %d: unexpected end of line", idx+1)
			}
			val, err := strconv.Atoi(fields[pos])
			pos++
			if err != nil {
				return 0, fmt.Errorf("line %d: parse int: %w", idx+1, err)
			}
			return val, nil
		}

		n, err := parseInt()
		if err != nil {
			return nil, err
		}
		if len(fields) < 1+n {
			return nil, fmt.Errorf("line %d: not enough parent values", idx+1)
		}
		parents := make([]int, n+1)
		for i := 1; i <= n; i++ {
			v, err := parseInt()
			if err != nil {
				return nil, fmt.Errorf("line %d: parent[%d]: %w", idx+1, i, err)
			}
			parents[i] = v
		}
		m, err := parseInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: m: %w", idx+1, err)
		}
		if len(fields) != pos+2*m {
			return nil, fmt.Errorf("line %d: expected %d query values, got %d", idx+1, 2*m, len(fields)-pos)
		}
		queries := make([][2]int, m)
		for i := 0; i < m; i++ {
			v, err := parseInt()
			if err != nil {
				return nil, fmt.Errorf("line %d: query %d v: %w", idx+1, i+1, err)
			}
			p, err := parseInt()
			if err != nil {
				return nil, fmt.Errorf("line %d: query %d p: %w", idx+1, i+1, err)
			}
			queries[i] = [2]int{v, p}
		}
		cases = append(cases, testCase{
			n:       n,
			parents: parents,
			queries: queries,
		})
	}
	return cases, nil
}

// solveCase mirrors the reference solution logic from 208E.go.
func solveCase(n int, parent []int, queries [][2]int) []int {
	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		p := parent[i]
		if p > 0 {
			children[p] = append(children[p], i)
		}
	}
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	depthList := make([][]int, n+2)
	t := 0
	var dfs func(int)
	dfs = func(u int) {
		t++
		tin[u] = t
		d := depth[u]
		if d >= len(depthList) {
			depthList = append(depthList, []int{})
		}
		depthList[d] = append(depthList[d], t)
		for _, v := range children[u] {
			depth[v] = d + 1
			dfs(v)
		}
		tout[u] = t
	}
	for i := 1; i <= n; i++ {
		if parent[i] == 0 {
			depth[i] = 0
			dfs(i)
		}
	}
	const LOG = 18
	up := make([][]int, LOG)
	up[0] = make([]int, n+1)
	for i := 1; i <= n; i++ {
		up[0][i] = parent[i]
	}
	for k := 1; k < LOG; k++ {
		up[k] = make([]int, n+1)
		for i := 1; i <= n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
	res := make([]int, len(queries))
	for qi, q := range queries {
		v, p := q[0], q[1]
		u := v
		for k := 0; k < LOG && u > 0; k++ {
			if (p>>k)&1 == 1 {
				u = up[k][u]
			}
		}
		if u == 0 {
			res[qi] = 0
			continue
		}
		D := depth[u] + p
		if D >= len(depthList) {
			res[qi] = 0
			continue
		}
		arr := depthList[D]
		l := sort.Search(len(arr), func(i int) bool { return arr[i] >= tin[u] })
		r := sort.Search(len(arr), func(i int) bool { return arr[i] > tout[u] })
		cnt := r - l
		if cnt > 0 {
			cnt--
		}
		res[qi] = cnt
	}
	return res
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(tc.parents[i]))
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(tc.queries))
	for _, q := range tc.queries {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		expect := solveCase(tc.n, tc.parents, tc.queries)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(buildInput(tc))
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		gotFields := strings.Fields(strings.TrimSpace(out.String()))
		if len(gotFields) != len(expect) {
			fmt.Printf("test %d: expected %d values got %d\n", idx+1, len(expect), len(gotFields))
			os.Exit(1)
		}
		for i, g := range gotFields {
			var val int
			if _, err := fmt.Sscan(g, &val); err != nil {
				fmt.Printf("test %d: parse output: %v\n", idx+1, err)
				os.Exit(1)
			}
			if val != expect[i] {
				fmt.Printf("test %d failed\nexpected: %v\ngot:      %v\n", idx+1, expect, gotFields)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
