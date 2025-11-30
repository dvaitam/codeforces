package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const maxN = 400000 + 5

const testcasesRaw = `
3 3 1 2 1 R
2 1 2 2 1 B 1 1 R
2 4 3 1 4 R 1 2 B 1 1 R
2 2 2 2 2 R 1 1 R
4 3 1 3 2 R
2 3 1 2 2 R
3 1 2 2 1 B 3 1 R
4 4 2 1 3 R 3 4 R
4 3 4 1 2 R 2 3 R 1 1 B 3 3 B
3 4 1 3 3 B
1 4 1 1 3 B
2 3 3 1 2 B 2 1 R 1 3 B
4 2 1 1 2 R
2 3 3 2 3 B 1 1 R 2 2 R
2 2 4 1 1 B 2 1 R 2 2 R 1 2 B
3 3 4 2 3 B 2 2 B 3 2 B 1 3 B
3 3 2 3 1 R 2 2 B
1 4 2 1 3 B 1 2 R
4 2 2 2 1 B 1 1 R
3 2 4 1 1 B 2 2 B 3 1 R 1 2 R
3 4 3 2 1 B 1 3 R 2 3 B
3 2 4 2 2 B 1 1 B 3 2 R 2 1 B
1 3 3 1 2 B 1 3 B 1 1 B
1 1 1 1 1 R
2 4 1 2 3 R
4 1 2 4 1 R 3 1 B
4 1 1 2 1 B
3 3 1 2 2 B
2 1 1 2 1 R
2 1 2 2 1 B 1 1 B
1 3 1 1 3 B
4 3 2 2 2 B 3 3 R
4 1 1 2 1 R
1 1 1 1 1 B
4 1 1 1 1 R
2 3 4 1 1 R 2 1 R 1 3 B 1 2 R
1 4 2 1 4 R 1 3 B
1 2 2 1 2 R 1 1 R
3 3 2 1 2 B 3 2 R
2 1 2 1 1 R 2 1 B
1 1 1 1 1 B
2 1 1 2 1 R
4 3 2 3 2 B 4 2 R
3 3 4 2 3 R 3 3 R 3 1 B 1 1 R
1 2 1 1 1 B
4 3 1 4 2 B
4 3 4 4 1 R 2 3 R 1 2 B 3 1 B
2 3 2 1 2 R 2 2 B
1 1 1 1 1 R
3 3 4 3 1 R 1 1 R 3 2 B 1 2 B
3 1 3 3 1 R 1 1 R 2 1 R
1 3 2 1 2 R 1 1 R
4 2 2 1 2 R 1 1 R
2 1 2 2 1 B 1 1 B
2 1 1 2 1 R
3 4 3 1 2 R 2 2 R 2 1 B
2 1 1 2 1 R
1 4 3 1 1 B 1 2 R 1 4 B
1 4 2 1 1 B 1 2 B
4 3 3 1 1 R 3 2 B 4 1 B
2 3 3 1 1 R 1 2 R 2 2 R
3 3 4 1 2 B 1 3 R 3 2 B 3 1 R
2 4 1 1 1 R
2 4 3 2 4 R 1 3 R 2 3 R
3 3 2 3 1 B 2 2 R
2 2 4 2 1 R 1 2 R 2 2 B 1 1 B
1 1 1 1 1 B
1 3 3 1 3 B 1 2 R 1 1 R
1 3 3 1 2 R 1 3 R 1 1 R
3 3 2 2 3 R 1 1 B
4 4 1 4 4 R
2 3 4 2 1 B 1 1 R 2 2 B 1 3 B
2 2 1 2 2 R
4 4 2 3 4 B 1 3 R
4 1 2 1 1 R 4 1 B
1 3 3 1 1 B 1 3 R 1 2 R
3 1 3 2 1 R 3 1 B 1 1 R
2 3 3 1 2 R 2 2 R 1 3 B
2 1 2 2 1 R 1 1 B
3 2 3 2 1 B 3 1 B 2 2 R
4 1 2 4 1 B 1 1 B
1 4 4 1 4 R 1 2 R 1 1 R 1 3 B
1 1 1 1 1 R
2 1 2 2 1 B 1 1 B
2 3 2 2 1 B 1 2 B
2 3 2 2 2 B 1 3 R
2 4 1 1 4 R
4 2 2 1 2 B 4 1 B
2 1 2 1 1 R 2 1 B
2 3 3 1 1 B 1 2 R 1 3 B
3 3 3 2 3 B 1 3 R 2 1 R
4 2 2 4 1 R 1 1 B
2 2 3 1 1 R 2 2 R 2 1 R
3 3 4 2 2 R 3 3 R 3 2 R 3 1 B
1 3 3 1 2 R 1 3 R 1 1 R
1 1 1 1 1 B
3 3 4 1 3 R 2 2 R 3 2 B 2 3 B
2 2 1 2 1 R
4 1 4 1 1 B 4 1 R 3 1 B 2 1 B
2 4 1 1 3 R
`

type testCase struct {
	n   int
	m   int
	q   int
	ops []op
}

type op struct {
	x int
	y int
	c string
}

type edge struct {
	x, y, t int
}

var (
	fa  [maxN]int
	sz  [maxN]int
	ans int64

	e   [maxN][]int
	dfn [maxN]int
	low [maxN]int
	st  [maxN]int
	tp  int
	co  [maxN]int
	tim int
	vis [maxN]bool

	outputs []int64
)

func gf(x int) int {
	if fa[x] == x {
		return x
	}
	fa[x] = gf(fa[x])
	return fa[x]
}

func merge(x, y int) {
	x = gf(x)
	y = gf(y)
	if x == y {
		return
	}
	if sz[x] > 1 {
		ans -= int64(sz[x]) * int64(sz[x])
	}
	if sz[y] > 1 {
		ans -= int64(sz[y]) * int64(sz[y])
	}
	fa[y] = x
	sz[x] += sz[y]
	ans += int64(sz[x]) * int64(sz[x])
}

func tarjan(x int) {
	tim++
	dfn[x] = tim
	low[x] = tim
	tp++
	st[tp] = x
	vis[x] = true
	for _, v := range e[x] {
		if dfn[v] == 0 {
			tarjan(v)
			if low[v] < low[x] {
				low[x] = low[v]
			}
		} else if vis[v] {
			if dfn[v] < low[x] {
				low[x] = dfn[v]
			}
		}
	}
	if low[x] == dfn[x] {
		co[x] = x
		for st[tp] != x {
			co[st[tp]] = x
			vis[st[tp]] = false
			tp--
		}
		tp--
		vis[x] = false
	}
}

func solveRec(l, r, q int, ed []edge) {
	if l == r {
		if l > q {
			return
		}
		for _, v := range ed {
			merge(v.x, v.y)
		}
		outputs = append(outputs, ans)
		return
	}
	mid := (l + r) >> 1
	tim = 0
	tp = 0
	var el, er []edge
	for i := range ed {
		v := &ed[i]
		v.x = gf(v.x)
		v.y = gf(v.y)
		e[v.x] = nil
		e[v.y] = nil
		dfn[v.x] = 0
		dfn[v.y] = 0
	}
	for _, v := range ed {
		if v.t <= mid {
			e[v.x] = append(e[v.x], v.y)
		}
	}
	for _, v := range ed {
		if v.t <= mid {
			if dfn[v.x] == 0 {
				tarjan(v.x)
			}
			if dfn[v.y] == 0 {
				tarjan(v.y)
			}
			if co[v.x] == co[v.y] {
				el = append(el, v)
				continue
			}
		}
		er = append(er, v)
	}
	solveRec(l, mid, q, el)
	solveRec(mid+1, r, q, er)
}

func solveCase(tc testCase) string {
	total := tc.n + tc.m
	for i := 1; i <= total; i++ {
		fa[i] = i
		sz[i] = 1
		e[i] = nil
		dfn[i] = 0
		low[i] = 0
		co[i] = 0
		vis[i] = false
	}
	ans = 0
	outputs = outputs[:0]
	edges := make([]edge, tc.q)
	for i, op := range tc.ops {
		if op.c == "R" {
			edges[i] = edge{x: op.y + tc.n, y: op.x, t: i + 1}
		} else {
			edges[i] = edge{x: op.x, y: op.y + tc.n, t: i + 1}
		}
	}
	solveRec(1, tc.q+1, tc.q, edges)
	res := make([]string, len(outputs))
	for i, v := range outputs {
		res[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(res, "\n")
}

func parseTests(raw string) ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("bad line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad n in line: %q", line)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("bad m in line: %q", line)
		}
		q, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("bad q in line: %q", line)
		}
		expectLen := 3 + 3*q
		if len(fields) != expectLen {
			return nil, fmt.Errorf("unexpected field count in line: %q", line)
		}
		ops := make([]op, q)
		idx := 3
		for i := 0; i < q; i++ {
			x, _ := strconv.Atoi(fields[idx])
			y, _ := strconv.Atoi(fields[idx+1])
			c := fields[idx+2]
			ops[i] = op{x: x, y: y, c: c}
			idx += 3
		}
		tests = append(tests, testCase{n: n, m: m, q: q, ops: ops})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.q))
	for _, op := range tc.ops {
		sb.WriteString(fmt.Sprintf("%d %d %s\n", op.x, op.y, op.c))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		want := solveCase(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
