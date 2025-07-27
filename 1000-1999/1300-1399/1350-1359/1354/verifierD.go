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

func buildOracle() (string, error) {
	exe := "oracleD.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1354D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

type testCase struct {
	n       int
	q       int
	a       []int
	queries []int
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(n) + 1
	}
	sort.Ints(a)
	cur := n
	queries := make([]int, q)
	for i := 0; i < q; i++ {
		if cur == 0 || rng.Intn(2) == 0 {
			queries[i] = rng.Intn(n+q) + 1
			cur++
		} else {
			idx := rng.Intn(cur) + 1
			queries[i] = -idx
			cur--
		}
	}
	return testCase{n, q, a, queries}
}

func deterministicCases() []testCase {
	return []testCase{
		{1, 1, []int{1}, []int{-1}},
		{2, 2, []int{1, 2}, []int{-1, -1}},
		{2, 2, []int{1, 1}, []int{3, -1}},
	}
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
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
		return fmt.Errorf("expected %s got %s", expected, result)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
