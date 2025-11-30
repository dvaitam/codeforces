package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt.
const testcasesDData = `
3 2
1 3 1
2 3 16

8 10
1 5 1
2 3 19
7 8 5
7 3 8
6 2 2
3 8 7
5 4 7
5 7 10
5 8 20
3 7 1

5 8
4 5 8
4 2 1
3 1 9
5 2 11
5 1 7
1 3 3
4 5 12
5 3 10

4 6
2 1 10
3 2 2
1 4 19
2 4 10
3 2 13
4 1 2

5 7
5 4 20
5 3 14
2 3 14
5 1 20
4 5 11
5 3 1
3 5 10

2 1
1 2 9

7 9
6 7 9
5 6 5
1 7 1
3 7 13
3 2 12
2 3 12
7 5 9
3 5 8
3 5 20

3 3
3 1 15
1 3 14
`

const inf64 = int64(4e18)

type edge struct {
	u, v int
	w    int64
}

type testCase struct {
	n     int
	m     int
	edges []edge
}

// parseTestcases reads embedded data into structured cases.
func parseTestcases() ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(testcasesDData))
	var cases []testCase
	for {
		line, ok := nextNonEmpty(sc)
		if !ok {
			break
		}
		var n, m int
		if _, err := fmt.Sscan(line, &n, &m); err != nil {
			return nil, fmt.Errorf("parse n m: %w", err)
		}
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("unexpected EOF reading edges")
			}
			ln := strings.TrimSpace(sc.Text())
			if ln == "" {
				i--
				continue
			}
			var u, v int
			var w int64
			if _, err := fmt.Sscan(ln, &u, &v, &w); err != nil {
				return nil, fmt.Errorf("parse edge %d: %w", i+1, err)
			}
			edges[i] = edge{u - 1, v - 1, w}
		}
		// consume possible blank line between cases
		nextNonEmpty(sc)
		cases = append(cases, testCase{n: n, m: m, edges: edges})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

// nextNonEmpty returns the next non-empty trimmed line.
func nextNonEmpty(sc *bufio.Scanner) (string, bool) {
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			return line, true
		}
	}
	return "", false
}

// solveEdge returns minimal maximum distance for nodes from a point on edge u-v of weight w.
func solveEdge(u, v int, w int64, dist [][]int64) float64 {
	n := len(dist)
	a := make([]float64, n)
	b := make([]float64, n)
	bw := float64(w)
	type event struct {
		t   float64
		idx int
	}
	events := make([]event, n)
	for i := 0; i < n; i++ {
		a[i] = float64(dist[i][u])
		b[i] = float64(dist[i][v])
		events[i] = event{(b[i] + bw - a[i]) / 2.0, i}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })

	moved := make([]bool, n)
	alpha := make([]float64, n)
	copy(alpha, a)
	beta := make([]float64, n)

	maxAlpha := func() float64 {
		res := -1e30
		for i := 0; i < n; i++ {
			if moved[i] {
				continue
			}
			if alpha[i] > res {
				res = alpha[i]
			}
		}
		return res
	}
	maxBeta := func() float64 {
		res := -1e30
		for i := 0; i < n; i++ {
			if beta[i] > res {
				res = beta[i]
			}
		}
		return res
	}

	k := 0
	for k < n && events[k].t <= 0 {
		i := events[k].idx
		moved[i] = true
		beta[i] = b[i] + bw
		k++
	}
	prevT := 0.0
	best := math.Inf(1)

	process := func(L, R float64) {
		if L > bw || R < 0 {
			return
		}
		l := L
		if l < 0 {
			l = 0
		}
		r := R
		if r > bw {
			r = bw
		}
		if l > r {
			return
		}
		aMax := maxAlpha()
		bMax := maxBeta()
		t := (bMax - aMax) / 2.0
		if t < l {
			t = l
		} else if t > r {
			t = r
		}
		d1 := aMax + t
		d2 := bMax - t
		if d2 > d1 {
			d1 = d2
		}
		if d1 < best {
			best = d1
		}
	}

	for k < n {
		currT := events[k].t
		process(prevT, currT)
		for k < n && events[k].t == currT {
			i := events[k].idx
			moved[i] = true
			beta[i] = b[i] + bw
			k++
		}
		prevT = currT
	}
	process(prevT, bw)
	return best
}

// solve mirrors 266D.go for one test case.
func solve(tc testCase) float64 {
	n, m := tc.n, tc.m
	edges := tc.edges
	dist := make([][]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = inf64
			}
		}
	}
	for _, e := range edges {
		if e.w < dist[e.u][e.v] {
			dist[e.u][e.v] = e.w
			dist[e.v][e.u] = e.w
		}
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == inf64 {
				continue
			}
			for j := 0; j < n; j++ {
				d := dist[i][k] + dist[k][j]
				if d < dist[i][j] {
					dist[i][j] = d
				}
			}
		}
	}
	best := math.Inf(1)
	for u := 0; u < n; u++ {
		mx := int64(0)
		for i := 0; i < n; i++ {
			if dist[u][i] > mx {
				mx = dist[u][i]
			}
		}
		d := float64(mx)
		if d < best {
			best = d
		}
	}
	for _, e := range edges {
		d2 := solveEdge(e.u, e.v, e.w, dist)
		if d2 < best {
			best = d2
		}
	}
	return best
}

func runCandidate(bin string, tc testCase) (float64, error) {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
	for _, e := range tc.edges {
		fmt.Fprintf(&input, "%d %d %d\n", e.u+1, e.v+1, e.w)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	valStr := strings.TrimSpace(out.String())
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parse output %q: %v", valStr, err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if math.Abs(expect-got) > 1e-6*math.Max(1.0, math.Abs(expect)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.10f got %.10f\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
