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

func buildOracle() (string, error) {
	exe := "oracleE.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1354E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

type testCase struct {
	n     int
	m     int
	n1    int
	n2    int
	n3    int
	edges [][2]int
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 3
	n1 := rng.Intn(n + 1)
	n2 := rng.Intn(n - n1 + 1)
	n3 := n - n1 - n2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	used := make(map[[2]int]bool)
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, key)
	}
	return testCase{n, m, n1, n2, n3, edges}
}

func deterministicCases() []testCase {
	return []testCase{
		{3, 0, 1, 1, 1, nil},
		{4, 2, 1, 1, 2, [][2]int{{1, 2}, {3, 4}}},
	}
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n1, tc.n2, tc.n3))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func runCase(bin, oracle string, tc testCase) error {
	input := tc.input()
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	result := strings.TrimSpace(out.String())
	if result != expected {
		return fmt.Errorf("expected:\n%s\n---\ngot:\n%s", expected, result)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := deterministicCases()
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, oracle, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
