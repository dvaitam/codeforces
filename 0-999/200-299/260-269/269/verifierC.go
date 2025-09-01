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
	// build tree to ensure connectivity
	edges := make([][3]int, 0, m)
	parent := rng.Perm(n - 1)
	for i := 2; i <= n; i++ {
		p := parent[i-2] + 1
		w := rng.Intn(20) + 1
		edges = append(edges, [3]int{p, i, w})
	}
	used := make(map[[2]int]bool)
	for _, e := range edges {
		u, v := e[0], e[1]
		if u > v {
			u, v = v, u
		}
		used[[2]int{u, v}] = true
	}
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
		w := rng.Intn(20) + 1
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
	sumFlow := make([]int64, n+1)
	for i := 0; i < m; i++ {
		u := atoi(fields[p])
		v := atoi(fields[p+1])
		c := atoi(fields[p+2])
		p += 3
		a[i], b[i], w[i] = u, v, int64(c)
		sumFlow[u] += int64(c)
		sumFlow[v] += int64(c)
	}
	// parse output: expect m numbers 0/1
	outs := strings.Fields(output)
	if len(outs) < m {
		return fmt.Errorf("missing orientations: need %d got %d", m, len(outs))
	}
	outSum := make([]int64, n+1)
	for i := 0; i < m; i++ {
		var d int
		if _, err := fmt.Sscan(outs[i], &d); err != nil {
			return fmt.Errorf("bad orientation at %d", i+1)
		}
		if d != 0 && d != 1 {
			return fmt.Errorf("orientation must be 0/1 at %d", i+1)
		}
		if d == 0 {
			outSum[a[i]] += w[i]
		} else {
			outSum[b[i]] += w[i]
		}
	}
	// verify balances
	if outSum[1] != sumFlow[1] {
		return fmt.Errorf("node 1 out=%d sum=%d", outSum[1], sumFlow[1])
	}
	if outSum[n] != 0 {
		return fmt.Errorf("node %d out must be 0, got %d", n, outSum[n])
	}
	for v := 2; v <= n-1; v++ {
		if outSum[v]*2 != sumFlow[v] {
			return fmt.Errorf("node %d out*2=%d sum=%d", v, outSum[v]*2, sumFlow[v])
		}
	}
	// ensure no extra significant tokens beyond m
	if len(outs) > m {
		// allow trailing whitespace and lines, but if extra tokens exist, error
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
