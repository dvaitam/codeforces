package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *dsu) union(a, b int) int {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return a
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
	return a
}

type testCase struct {
	n, m    int
	types   []int
	queries [][2]int
}

const testcaseData = `
100
4 3
1 3 2 2
1 2
1 3
3 3
1 1 3
1 2
1 3
4 4
1 3 4 2
1 2
1 3
1 4
3 3
3 2 2
1 2
1 3
5 4
3 1 3 4 2
1 2
1 3
1 4
3 3
2 2 3
1 2
1 3
4 4
3 3 4 1
1 2
1 3
1 4
5 5
5 2 3 3 1
1 2
1 3
1 4
1 5
5 5
5 1 4 3 2
1 2
1 3
1 4
1 5
5 2
1 2 2 1 2
1 2
5 5
5 4 3 1 1
1 2
1 3
1 4
1 5
4 4
1 3 2 4
1 2
1 3
1 4
2 2
2 1
1 2
5 4
1 3 3 4 2
1 2
1 3
1 4
4 2
1 2 2 1
1 2
5 4
2 3 1 4 1
1 2
1 3
1 4
4 3
2 2 3 1
1 2
1 3
5 5
3 4 1 1 5
1 2
1 3
1 4
1 5
5 3
1 3 2 3 1
1 2
1 3
4 4
1 3 4 2
1 2
1 3
1 4
5 5
2 3 4 5 1
1 2
1 3
1 4
1 5
5 5
3 2 5 2 1
1 2
1 3
1 4
1 5
4 4
2 1 4 3
1 2
1 3
1 4
4 4
3 1 4 2
1 2
1 3
1 4
5 5
4 2 1 5 3
1 2
1 3
1 4
1 5
5 4
1 3 2 4 4
1 2
1 3
1 4
3 3
2 3 1
1 2
1 3
4 4
2 3 4 1
1 2
1 3
1 4
3 2
1 1 2
1 2
5 2
2 1 1 1 1
1 2
5 4
4 1 1 2 3
1 2
1 3
1 4
4 4
1 3 4 1
1 2
1 3
1 4
2 2
1 2
1 2
5 4
3 4 4 1 1
1 2
1 3
1 4
5 4
4 4 1 2 3
1 2
1 3
1 4
5 5
1 5 2 3 2
1 2
1 3
1 4
1 5
5 4
3 3 4 2 1
1 2
1 3
1 4
2 2
1 2
1 2
4 3
3 2 1 1
1 2
1 3
4 2
2 2 2 1
1 2
4 3
3 1 3 2
1 2
1 3
5 5
5 4 1 2 3
1 2
1 3
1 4
1 5
4 4
1 4 2 3
1 2
1 3
1 4
5 4
2 3 4 1 1
1 2
1 3
1 4
3 2
1 2 1
1 2
3 2
2 1 2
1 2
5 5
4 2 1 3 5
1 2
1 3
1 4
1 5
5 5
5 4 3 1 2
1 2
1 3
1 4
1 5
4 2
1 2 2 2
1 2
5 5
2 5 1 4 3
1 2
1 3
1 4
1 5
4 4
1 3 2 4
1 2
1 3
1 4
3 3
3 1 1
1 2
1 3
5 5
4 3 5 2 1
1 2
1 3
1 4
1 5
2 2
2 1
1 2
4 4
2 4 3 1
1 2
1 3
1 4
3 2
2 1 1
1 2
5 4
1 4 2 2 3
1 2
1 3
1 4
5 5
4 3 2 5 1
1 2
1 3
1 4
1 5
4 2
1 2 1 1
1 2
3 3
2 1 3
1 2
1 3
5 2
1 2 2 2 1
1 2
3 3
2 1 3
1 2
1 3
5 5
5 3 4 1 3
1 2
1 3
1 4
1 5
4 3
1 3 2 1
1 2
1 3
5 5
1 2 1 4 5
1 2
1 3
1 4
1 5
3 3
1 2 3
1 2
1 3
5 3
2 1 3 3 1
1 2
1 3
3 2
2 1 1
1 2
5 2
1 1 2 1 2
1 2
5 4
3 1 2 4 4
1 2
1 3
1 4
5 4
1 4 3 2 2
1 2
1 3
1 4
5 5
3 2 5 4 1
1 2
1 3
1 4
1 5
4 4
1 2 4 3
1 2
1 3
1 4
3 2
1 1 2
1 2
5 3
2 1 3 2 1
1 2
1 3
5 5
3 1 4 5 1
1 2
1 3
1 4
1 5
4 4
1 3 4 2
1 2
1 3
1 4
4 2
2 2 1 1
1 2
3 3
2 1 3
1 2
1 3
5 4
1 2 2 4 3
1 2
1 3
1 4
5 5
3 5 4 1 1
1 2
1 3
1 4
1 5
4 4
3 2 1 4
1 2
1 3
1 4
4 4
4 1 3 2
1 2
1 3
1 4
2 2
1 2
1 2
5 4
1 3 4 2 2
1 2
1 3
1 4
5 5
3 3 1 5 2
1 2
1 3
1 4
1 5
5 2
1 1 2 2 1
1 2
5 2
1 1 1 1 2
1 2
2 2
1 2
1 2
5 3
1 1 1 3 2
1 2
1 3
4 3
3 2 1 3
1 2
1 3
5 3
2 3 1 3 3
1 2
1 3
5 4
1 2 3 2 4
1 2
1 3
1 4
5 5
1 4 2 2 5
1 2
1 3
1 4
1 5
4 2
1 1 2 1
1 2
4 2
2 2 1 2
1 2
4 4
2 3 4 1
1 2
1 3
1 4
4 2
2 1 1 2
1 2
5 2
1 2 1 1 1
1 2
5 3
3 1 2 1 1
1 2
1 3
`

func parseTestcases() ([]testCase, error) {
	tokens := strings.Fields(testcaseData)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	idx := 0
	next := func() (string, error) {
		if idx >= len(tokens) {
			return "", fmt.Errorf("unexpected end of data")
		}
		val := tokens[idx]
		idx++
		return val, nil
	}

	tStr, err := next()
	if err != nil {
		return nil, err
	}
	t, err := strconv.Atoi(tStr)
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}

	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		nStr, err := next()
		if err != nil {
			return nil, fmt.Errorf("case %d: missing n", caseIdx+1)
		}
		mStr, err := next()
		if err != nil {
			return nil, fmt.Errorf("case %d: missing m", caseIdx+1)
		}
		n, err := strconv.Atoi(nStr)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %v", caseIdx+1, err)
		}
		m, err := strconv.Atoi(mStr)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad m: %v", caseIdx+1, err)
		}
		types := make([]int, n+1)
		for i := 1; i <= n; i++ {
			vStr, err := next()
			if err != nil {
				return nil, fmt.Errorf("case %d: missing type %d", caseIdx+1, i)
			}
			v, err := strconv.Atoi(vStr)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad type %d: %v", caseIdx+1, i, err)
			}
			types[i] = v
		}
		queries := make([][2]int, m-1)
		for i := 0; i < m-1; i++ {
			aStr, err := next()
			if err != nil {
				return nil, fmt.Errorf("case %d: missing query a", caseIdx+1)
			}
			bStr, err := next()
			if err != nil {
				return nil, fmt.Errorf("case %d: missing query b", caseIdx+1)
			}
			a, err := strconv.Atoi(aStr)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query a: %v", caseIdx+1, err)
			}
			b, err := strconv.Atoi(bStr)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad query b: %v", caseIdx+1, err)
			}
			queries[i] = [2]int{a, b}
		}
		cases = append(cases, testCase{n: n, m: m, types: types, queries: queries})
	}
	return cases, nil
}

func solveCase(tc testCase) []int {
	n, m := tc.n, tc.m
	t := tc.types
	items := make([][]int, m+1)
	for i := 1; i <= n; i++ {
		items[t[i]] = append(items[t[i]], i)
	}

	dsu := newDSU(m)
	for i := 1; i <= m; i++ {
		dsu.size[i] = len(items[i])
	}

	ans := 0
	for i := 1; i < n; i++ {
		if t[i] != t[i+1] {
			ans++
		}
	}

	res := make([]int, 0, m)
	res = append(res, ans)

	for _, q := range tc.queries {
		a := dsu.find(q[0])
		b := dsu.find(q[1])
		if a != b {
			if dsu.size[a] < dsu.size[b] {
				a, b = b, a
			}
			for _, x := range items[b] {
				if x > 1 && dsu.find(t[x-1]) == a {
					ans--
				}
				if x < n && dsu.find(t[x+1]) == a {
					ans--
				}
			}
			dsu.parent[b] = a
			dsu.size[a] += dsu.size[b]
			items[a] = append(items[a], items[b]...)
			items[b] = nil
		}
		res = append(res, ans)
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for j := 1; j <= tc.n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.types[j]))
		}
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
		}

		expected := solveCase(tc)
		gotStr, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(strings.TrimSpace(gotStr))
		if len(gotFields) != len(expected) {
			fmt.Printf("case %d failed: expected %d numbers got %d\n", i+1, len(expected), len(gotFields))
			os.Exit(1)
		}
		for idx, exp := range expected {
			val, err := strconv.Atoi(gotFields[idx])
			if err != nil || val != exp {
				fmt.Printf("case %d failed at position %d: expected %d got %s\n", i+1, idx+1, exp, gotFields[idx])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
