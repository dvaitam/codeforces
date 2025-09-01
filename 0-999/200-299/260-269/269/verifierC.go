package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
	// build a proper tree to ensure connectivity (no self-loops)
	edges := make([][3]int, 0, m)
	used := make(map[[2]int]bool)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1 // connect to a previous node, guarantees p != i
		u, v := p, i
		if u > v {
			u, v = v, u
		}
		if !used[[2]int{u, v}] {
			used[[2]int{u, v}] = true
			w := 2 * (rng.Intn(10) + 1) // even weight
			edges = append(edges, [3]int{u, v, w})
		}
	}
	// add extra random edges without duplicates or self-loops
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		if used[[2]int{u, v}] {
			continue
		}
		used[[2]int{u, v}] = true
		w := 2 * (rng.Intn(10) + 1) // even weight
		edges = append(edges, [3]int{u, v, w})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return testCase{input: sb.String()}
}

func run(bin, input string) (string, error) {
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

func verify(input, output string) error {
	// parse input
	fields := strings.Fields(input)
	p := 0
	if len(fields) < 2 {
		return fmt.Errorf("bad input")
	}
	n := atoi(fields[p])
	p++
	m := atoi(fields[p])
	p++
	a := make([]int, m)
	b := make([]int, m)
	w := make([]int64, m)
	for i := 0; i < m; i++ {
		u := atoi(fields[p])
		v := atoi(fields[p+1])
		c := atoi(fields[p+2])
		p += 3
		a[i], b[i], w[i] = u, v, int64(c)
	}
	// parse output: expect m numbers 0/1
	outs := strings.Fields(output)
	if len(outs) < m {
		return fmt.Errorf("missing orientations: need %d got %d", m, len(outs))
	}
	outSum := make([]int64, n+1) // total flow going out of each node
	inSum := make([]int64, n+1)  // total flow coming into each node
	for i := 0; i < m; i++ {
		var d int
		if _, err := fmt.Sscan(outs[i], &d); err != nil {
			return fmt.Errorf("bad orientation at %d", i+1)
		}
		if d != 0 && d != 1 {
			return fmt.Errorf("orientation must be 0/1 at %d", i+1)
		}
		u, v, ww := a[i], b[i], w[i]
		if u == v {
			// self-loop contributes equally to in and out; not present in CF input, but be robust
			outSum[u] += ww
			inSum[u] += ww
			continue
		}
		if d == 0 { // u -> v
			outSum[u] += ww
			inSum[v] += ww
		} else { // v -> u
			outSum[v] += ww
			inSum[u] += ww
		}
	}
	// internal nodes must satisfy in == out
	for v := 2; v <= n-1; v++ {
		if inSum[v] != outSum[v] {
			return fmt.Errorf("node %d out=%d in=%d", v, outSum[v], inSum[v])
		}
	}
	// ensure no extra tokens beyond m
	if len(outs) > m {
		return fmt.Errorf("extra output tokens: %v", outs[m:])
	}
	return nil
}

func atoi(s string) int {
	sign := 1
	i := 0
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		i = 1
	}
	v := 0
	for ; i < len(s); i++ {
		v = v*10 + int(s[i]-'0')
	}
	return v * sign
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 101)
	for i := 0; i < 101; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := verify(tc.input, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
