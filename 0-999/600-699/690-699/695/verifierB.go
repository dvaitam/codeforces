package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u, v int
	w    int
}

type testCase struct {
	n, m int
	s, t int
	e    []edge
}

func parseInput(input string) (testCase, error) {
	fields := strings.Fields(input)
	vals := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return testCase{}, fmt.Errorf("invalid integer %q", f)
		}
		vals[i] = v
	}
	if len(vals) < 4 {
		return testCase{}, fmt.Errorf("too few values in input")
	}
	tc := testCase{}
	tc.n, tc.m, tc.s, tc.t = vals[0], vals[1], vals[2], vals[3]
	need := 4 + 3*tc.m
	if len(vals) != need {
		return testCase{}, fmt.Errorf("bad input length: got %d values, need %d", len(vals), need)
	}
	tc.e = make([]edge, tc.m)
	idx := 4
	for i := 0; i < tc.m; i++ {
		tc.e[i] = edge{u: vals[idx], v: vals[idx+1], w: vals[idx+2]}
		idx += 3
	}
	return tc, nil
}

func disconnected(tc testCase, removed map[int]bool) bool {
	vis := make([]bool, tc.n+1)
	q := []int{tc.s}
	vis[tc.s] = true
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for i, e := range tc.e {
			id := i + 1
			if removed[id] {
				continue
			}
			to := 0
			if e.u == v {
				to = e.v
			} else if e.v == v {
				to = e.u
			}
			if to != 0 && !vis[to] {
				vis[to] = true
				q = append(q, to)
			}
		}
	}
	return !vis[tc.t]
}

func optimalCost(tc testCase) int {
	const inf = int(1e9)
	best := inf
	if disconnected(tc, map[int]bool{}) {
		best = 0
	}
	for i := 1; i <= tc.m; i++ {
		if disconnected(tc, map[int]bool{i: true}) && tc.e[i-1].w < best {
			best = tc.e[i-1].w
		}
	}
	for i := 1; i <= tc.m; i++ {
		for j := i + 1; j <= tc.m; j++ {
			cost := tc.e[i-1].w + tc.e[j-1].w
			if cost >= best {
				continue
			}
			if disconnected(tc, map[int]bool{i: true, j: true}) {
				best = cost
			}
		}
	}
	if best == inf {
		return -1
	}
	return best
}

func validateOutput(input, out string) error {
	tc, err := parseInput(input)
	if err != nil {
		return err
	}
	best := optimalCost(tc)

	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	vals := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("output is not integer-only: %q", f)
		}
		vals[i] = v
	}

	if vals[0] == -1 {
		if len(vals) != 1 {
			return fmt.Errorf("-1 output must contain only one integer")
		}
		if best != -1 {
			return fmt.Errorf("expected feasible answer with cost %d, got -1", best)
		}
		return nil
	}

	if best == -1 {
		return fmt.Errorf("expected -1, got feasible output")
	}
	if len(vals) < 2 {
		return fmt.Errorf("output must contain cost and number of edges")
	}
	printedCost, k := vals[0], vals[1]
	if k < 0 || k > 2 {
		return fmt.Errorf("invalid number of edges: %d", k)
	}
	if len(vals) != 2+k {
		return fmt.Errorf("edge list length mismatch: k=%d but got %d ids", k, len(vals)-2)
	}
	removed := make(map[int]bool)
	sum := 0
	for _, id := range vals[2:] {
		if id < 1 || id > tc.m {
			return fmt.Errorf("edge id out of range: %d", id)
		}
		if removed[id] {
			return fmt.Errorf("duplicate edge id: %d", id)
		}
		removed[id] = true
		sum += tc.e[id-1].w
	}
	if sum != printedCost {
		return fmt.Errorf("printed cost %d does not match chosen edges cost %d", printedCost, sum)
	}
	if !disconnected(tc, removed) {
		return fmt.Errorf("chosen edges do not disconnect %d and %d", tc.s, tc.t)
	}
	if printedCost != best {
		return fmt.Errorf("expected minimal cost %d got %d", best, printedCost)
	}
	return nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := n - 1 + rng.Intn(maxEdges-(n-1)+1)
	s := rng.Intn(n) + 1
	t := rng.Intn(n-1) + 1
	if t >= s {
		t++
	}
	type pair struct{ u, v int }
	edges := make([][3]int, 0, m)
	used := make(map[pair]bool)
	for i := 2; i <= n; i++ {
		u := i
		v := rng.Intn(i-1) + 1
		w := rng.Intn(20) + 1
		edges = append(edges, [3]int{u, v, w})
		if u > v {
			u, v = v, u
		}
		used[pair{u, v}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		du, dv := u, v
		if du > dv {
			du, dv = dv, du
		}
		if used[pair{du, dv}] {
			continue
		}
		w := rng.Intn(20) + 1
		edges = append(edges, [3]int{u, v, w})
		used[pair{du, dv}] = true
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	return sb.String()
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if err := validateOutput(input, got); err != nil {
		return fmt.Errorf("%v\ngot:\n%s", err, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
