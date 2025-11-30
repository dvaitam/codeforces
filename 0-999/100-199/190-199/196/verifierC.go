package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
100
2
1 2
2 11
5 9
4
1 2
1 3
3 4
5 13
20 12
16 11
17 14
4
1 2
1 3
2 4
14 10
12 13
16 5
17 5
3
1 2
1 3
5 10
5 4
16 16
4
1 2
2 3
2 4
16 11
18 11
11 14
5 12
5
1 2
2 3
2 4
4 5
16 16
11 14
14 11
18 17
7 10
3
1 2
2 3
9 9
16 17
16 16
5
1 2
1 3
2 4
3 5
19 2
10 0
6 3
1 18
20 1
4
1 2
1 3
3 4
4 8
7 6
1 13
1 1
4
1 2
1 3
1 4
0 2
3 2
0 1
8 4
3
1 2
1 3
12 18
1 7
4 1
2
1 2
19 20
3 9
4
1 2
1 3
2 4
14 17
19 1
8 12
19 4
5
1 2
1 3
3 4
3 5
3 0
14 4
16 18
12 15
16 10
3
1 2
2 3
8 19
13 20
0 17
3
1 2
2 3
1 4
5 5
3 14
3
1 2
1 3
7 14
2 8
2 18
3
1 2
2 3
13 8
16 0
4 1
5
1 2
1 3
1 4
1 5
7 3
3 0
5 7
3 6
0 16
5
1 2
2 3
3 4
4 5
6 6
13 13
16 0
1 13
16 18
3
1 2
2 3
11 0
16 3
19 11
4
1 2
2 3
1 4
13 3
3 9
6 0
14 1
5
1 2
2 3
1 4
1 5
0 9
0 11
9 2
7 15
6 3
4
1 2
2 3
1 4
11 12
3 8
3 3
2 19
4
1 2
1 3
3 4
3 0
19 15
1 15
9 11
5
1 2
2 3
2 4
4 5
16 15
13 15
9 12
7 5
15 19
4
1 2
1 3
3 4
18 3
2 11
5 17
4 13
2
1 2
20 1
4 9
5
1 2
2 3
2 4
2 5
16 9
3 4
17 13
3 10
16 7
4
1 2
1 3
2 4
7 12
11 18
4 14
14 0
5
1 2
2 3
3 4
1 5
15 8
12 8
13 20
15 11
17 10
2
1 2
17 19
6 12
5
1 2
2 3
2 4
4 5
20 5
3 0
12 6
18 19
3 12
3
1 2
1 3
15 19
4 0
19 13
5
1 2
1 3
2 4
2 5
2 11
0 15
17 2
18 15
10 14
4
1 2
1 3
1 4
19 11
5 12
8 20
4 1
3
1 2
2 3
14 9
4 0
9 17
5
1 2
2 3
1 4
4 5
18 14
6 9
15 20
4 15
17 9
2
1 2
10 9
10 20
4
1 2
1 3
3 4
20 6
12 19
16 4
16 20
2
1 2
1 7
14 17
3
1 2
1 3
3 3
12 11
6 10
4
1 2
2 3
2 4
11 5
15 14
9 14
20 6
4
1 2
1 3
1 4
7 15
6 11
5 11
4 4
3
1 2
2 3
12 10
8 19
16 18
4
1 2
2 3
3 4
19 20
2 11
9 12
15 5
4
1 2
2 3
2 4
2 5
10 12
4 0
3 11
3
1 2
1 3
20 13
0 17
10 7
5
1 2
2 3
3 4
2 5
11 10
6 15
3 4
6 10
8 4
5
1 2
2 3
1 4
3 5
6 7
7 19
1 10
11 20
19 1
3
1 2
1 3
13 14
8 4
10 16
2
1 2
20 19
12 7
2
1 2
15 15
19 10
2
1 2
12 14
5 13
5
1 2
1 3
1 4
4 5
18 4
3 16
5 2
12 9
14 0
4
1 2
2 3
1 4
5 0
4 13
2 10
20 14
2
1 2
7 2
15 4
2
1 2
16 17
1 1
3
1 2
2 3
16 7
4 11
15 0
3
1 2
1 3
3 14
6 1
19 6
5
1 2
2 3
3 4
2 5
16 3
4 20
6 5
12 6
9 10
5
1 2
2 3
1 4
4 5
10 9
3 17
3 15
8 9
16 15
4
1 2
2 3
3 4
4 17
3 0
19 17
6 6
3
1 2
1 3
20 4
20 0
8 15
2
1 2
4 19
10 1
3
1 2
1 3
20 17
5 2
14 20
4
1 2
1 3
2 4
8 16
18 2
13 13
1 14
4
1 2
2 3
1 4
6 13
10 8
17 12
18 16
3
1 2
1 3
5 14
14 11
12 15
4
1 2
2 3
2 4
6 15
18 10
9 2
5 11
5
1 2
1 3
3 4
3 5
6 17
9 3
0 0
6 10
1 10
4
1 2
2 3
1 4
13 15
0 9
18 18
4 6
3
1 2
2 3
2 20
18 14
8 20
2
1 2
15 7
4 18
4
1 2
1 3
3 4
10 18
19 12
16 13
6 17
2
1 2
7 4
19 12
5
1 2
2 3
2 4
4 5
15 12
9 6
7 7
1 17
16 2
2
1 2
12 13
12 7
4
1 2
2 3
3 4
11 16
15 18
2 14
7 8
2
1 2
15 1
4 20
3
1 2
2 3
7 17
1 19
4 20
4
1 2
1 3
3 4
4 13
4 1
9 16
8 15
2
1 2
10 3
19 11
2
1 2
11 20
8 15
4
1 2
1 3
1 4
10 13
20 0
11 17
1 2
5
1 2
2 3
1 4
2 5
5 19
1 0
18 11
5 9
0 1
3
1 2
2 3
2 11
3 19
2 7
3
1 2
1 3
0 12
2 16
8 18
3
1 2
2 3
13 4
4 13
4 14
4
1 2
1 3
3 4
14 13
19 20
14 5
15 19
3
1 2
1 3
0 8
5 4
20 13
4
1 2
2 3
2 4
6 13
13 8
7 11
20 1
5
1 2
2 3
2 4
1 5
17 15
18 8
8 7
14 14
11 16
5
1 2
1 3
2 4
3 5
11 13
3 16
7 20
12 3
13 19
5
1 2
1 3
3 4
4 5
14 19
11 17
11 5
4 7
20 5
5
1 2
2 3
3 4
2 5
13 8
10 18
12 9
20 8
10 0
5
1 2
1 3
2 4
1 5
3 0
11 10
19 10
12 5
10 2
5
1 2
1 3
2 4
1 5
19 0
11 0
9 15
4 1
0 10
5
1 2
1 3
2 4
2 5
7 19
10 5
10 9
12 18
19 15
5
1 2
1 3
3 4
2 5
18 11
7 11
11 5
7 17
20 19
`

type point struct{ x, y int }

type testCase struct {
	n     int
	edges [][2]int
	pts   []point
}

func parseTests(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	idx := 0
	t, err := strconv.Atoi(fields[idx])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for c := 0; c < t; c++ {
		if idx >= len(fields) {
			return nil, fmt.Errorf("truncated at case %d", c+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("bad n in case %d", c+1)
		}
		idx++
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if idx+1 >= len(fields) {
				return nil, fmt.Errorf("truncated edges in case %d", c+1)
			}
			u, _ := strconv.Atoi(fields[idx])
			v, _ := strconv.Atoi(fields[idx+1])
			idx += 2
			edges[i] = [2]int{u, v}
		}
		pts := make([]point, n)
		for i := 0; i < n; i++ {
			if idx+1 >= len(fields) {
				return nil, fmt.Errorf("truncated points in case %d", c+1)
			}
			x, _ := strconv.Atoi(fields[idx])
			y, _ := strconv.Atoi(fields[idx+1])
			idx += 2
			pts[i] = point{x: x, y: y}
		}
		cases = append(cases, testCase{n: n, edges: edges, pts: pts})
	}
	return cases, nil
}

func vect(ax, ay, bx, by int) int64 {
	return int64(ax)*int64(by) - int64(ay)*int64(bx)
}

func solveCase(tc testCase) []int {
	n := tc.n
	g := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	type node struct{ x, y, order, idx int }
	mm := make([]node, n)
	for i, p := range tc.pts {
		mm[i] = node{x: p.x, y: p.y, order: 0, idx: i}
	}
	m := 0
	for i := 0; i < n; i++ {
		if mm[i].y > mm[m].y {
			m = i
		}
	}
	u := make([]bool, n+1)
	mc := make([]int, n+1)
	var dfs func(int) int
	dfs = func(v int) int {
		u[v] = true
		r := 1
		for _, w := range g[v] {
			if !u[w] {
				r += dfs(w)
			}
		}
		mc[v] = r
		return r
	}
	dfs(1)
	for i := range u {
		u[i] = false
	}
	px, py := 0, 0
	swap := func(i, j int) { mm[i], mm[j] = mm[j], mm[i] }
	var quickSort func(l, r int)
	quickSort = func(l, r int) {
		i, j := l, r
		x, y := mm[(i+j)/2].x, mm[(i+j)/2].y
		for i <= j {
			for i <= r && vect(mm[i].x-px, mm[i].y-py, x-px, y-py) < 0 {
				i++
			}
			for j >= l && vect(mm[j].x-px, mm[j].y-py, x-px, y-py) > 0 {
				j--
			}
			if i <= j {
				swap(i, j)
				i++
				j--
			}
		}
		if l < j {
			quickSort(l, j)
		}
		if i < r {
			quickSort(i, r)
		}
	}
	L := 0
	var rec func(a, v int)
	rec = func(a, v int) {
		mm[a].order = v
		u[v] = true
		swap(L, a)
		px = mm[L].x
		py = mm[L].y
		if mc[v] > 1 {
			quickSort(L+1, L+mc[v]-1)
		}
		L++
		for _, w := range g[v] {
			if !u[w] {
				rec(L, w)
			}
		}
	}
	rec(m, 1)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[mm[i].idx] = mm[i].order
	}
	return res
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(strconv.Itoa(e[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[1]))
		sb.WriteByte('\n')
	}
	for _, p := range tc.pts {
		sb.WriteString(strconv.Itoa(p.x))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(p.y))
		sb.WriteByte('\n')
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expectedArr := solveCase(tc)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotFields := strings.Fields(gotStr)
		if len(gotFields) != tc.n {
			fmt.Printf("test %d failed: expected %d numbers got %d\ninput:\n%s", idx+1, tc.n, len(gotFields), input)
			os.Exit(1)
		}
		for i, gf := range gotFields {
			val, err := strconv.Atoi(gf)
			if err != nil {
				fmt.Printf("test %d failed: bad int output\n", idx+1)
				os.Exit(1)
			}
			if val != expectedArr[i] {
				fmt.Printf("test %d failed: expected %d got %d\ninput:\n%s", idx+1, expectedArr[i], val, input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
