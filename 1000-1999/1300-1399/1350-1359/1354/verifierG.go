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
	exe := "oracleG.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1354G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

type testCase struct {
	n int
	k int
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 2
	k := rng.Intn(n/2) + 1
	return testCase{n, k}
}

func deterministicCases() []testCase {
	return []testCase{{2, 1}, {3, 1}}
}

func (tc testCase) input() string {
	return fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
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
		return fmt.Errorf("expected %q got %q", expected, result)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
