package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	in  string
	out string
}

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
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func solveCase(n int, p []int, pairs [][2]int) string {
	d := newDSU(n)
	for _, pr := range pairs {
		d.union(pr[0], pr[1])
	}
	compPos := make(map[int][]int)
	compVal := make(map[int][]int)
	for i := 1; i <= n; i++ {
		r := d.find(i)
		compPos[r] = append(compPos[r], i)
		compVal[r] = append(compVal[r], p[i-1])
	}
	ans := make([]int, n+1)
	for r := range compPos {
		pos := compPos[r]
		val := compVal[r]
		sort.Ints(pos)
		sort.Slice(val, func(i, j int) bool { return val[i] > val[j] })
		for i := 0; i < len(pos); i++ {
			ans[pos[i]] = val[i]
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", ans[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildCase(n int, p []int, pairs [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(pairs)))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, pr := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", pr[0], pr[1]))
	}
	return testCase{in: sb.String(), out: solveCase(n, p, pairs)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 2
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	maxPairs := n * (n - 1) / 2
	m := rng.Intn(maxPairs + 1)
	pairs := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	for len(pairs) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if used[key] {
			continue
		}
		used[key] = true
		pairs = append(pairs, key)
	}
	return buildCase(n, p, pairs)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.out)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase(1, []int{1}, nil))
	cases = append(cases, buildCase(2, []int{2, 1}, [][2]int{{1, 2}}))
	cases = append(cases, buildCase(3, []int{1, 2, 3}, nil))
	cases = append(cases, buildCase(3, []int{3, 2, 1}, [][2]int{{1, 2}, {2, 3}}))
	cases = append(cases, buildCase(4, []int{4, 3, 2, 1}, [][2]int{{1, 4}, {2, 3}}))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
