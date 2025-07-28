package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testE struct {
	n int
	s string
	a string
}

func buildOracle() (string, error) {
	exe := "oracleE.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1571E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func randBracket(rng *rand.Rand, l int) string {
	b := make([]byte, l)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = '('
		} else {
			b[i] = ')'
		}
	}
	return string(b)
}

func randBinary(rng *rand.Rand, l int) string {
	b := make([]byte, l)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func genTests() []testE {
	rng := rand.New(rand.NewSource(1571))
	tests := make([]testE, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(7) + 4
		s := randBracket(rng, n)
		a := randBinary(rng, n-3)
		tests = append(tests, testE{n, s, a})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", tc.n, tc.s, tc.a)
		want, err := runProg(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
