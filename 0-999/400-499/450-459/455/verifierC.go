package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "455C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBin(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

type dsu struct {
	p []int
	s []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &dsu{p: p, s: s}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(x, y int) bool {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return false
	}
	if d.s[rx] < d.s[ry] {
		rx, ry = ry, rx
	}
	d.p[ry] = rx
	d.s[rx] += d.s[ry]
	return true
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(n)
	q := rng.Intn(5) + 1
	d := newDSU(n)
	edges := [][2]int{}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if !d.union(u, v) {
			continue
		}
		edges = append(edges, [2]int{u, v})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			x := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("1 %d\n", x))
		} else {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d\n", x, y))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := make([]string, 0, 102)
	cases = append(cases, generateCase(rng))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, input := range cases {
		want, err := runBin(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", idx+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
