package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DSU struct {
	p    []int
	diff []int64
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	diff := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		diff[i] = 0
	}
	return &DSU{p: p, diff: diff}
}

func (d *DSU) find(x int) (int, int64) {
	if d.p[x] == x {
		return x, 0
	}
	r, dr := d.find(d.p[x])
	d.diff[x] += dr
	d.p[x] = r
	return r, d.diff[x]
}

func (d *DSU) unite(f, t int, delta int64) bool {
	rf, df := d.find(f)
	rt, dt := d.find(t)
	if rf == rt {
		return df-dt == delta
	}
	d.p[rf] = rt
	d.diff[rf] = delta + dt - df
	return true
}

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var n, m int
	fmt.Fscan(rdr, &n, &m)
	dsu := NewDSU(n)
	bad := 0
	for i := 1; i <= m; i++ {
		var f, t int
		var w, b int64
		fmt.Fscan(rdr, &f, &t, &w, &b)
		if bad != 0 {
			continue
		}
		delta := 2 * w * b
		if !dsu.unite(f, t, delta) {
			bad = i
		}
	}
	if bad != 0 {
		return fmt.Sprintf("BAD %d", bad)
	}
	r1, d1 := dsu.find(1)
	rn, dn := dsu.find(n)
	if r1 != rn {
		return "UNKNOWN"
	}
	diff := d1 - dn
	var ans int64
	if diff >= 0 {
		ans = (diff + 1) / 2
	} else {
		ans = (diff - 1) / 2
	}
	return fmt.Sprintf("%d", ans)
}

func genRandomCase() string {
	n := rand.Intn(4) + 2
	m := rand.Intn(5) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		f := rand.Intn(n) + 1
		t := rand.Intn(n) + 1
		w := rand.Intn(5) + 1
		b := rand.Intn(5) + 1
		fmt.Fprintf(&buf, "%d %d %d %d\n", f, t, w, b)
	}
	return buf.String()
}

func generateCases() []testCase {
	rand.Seed(5)
	cases := []testCase{}
	fixed := []string{
		"2 1\n1 2 1 1\n",
		"3 2\n1 2 1 1\n2 3 1 1\n",
	}
	for _, f := range fixed {
		cases = append(cases, testCase{f, compute(f)})
	}
	for len(cases) < 100 {
		inp := genRandomCase()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:%sexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
