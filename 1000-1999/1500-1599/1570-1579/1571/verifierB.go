package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testB struct {
	n  int
	a  int
	va int
	c  int
	vc int
	b  int
}

func buildOracle() (string, error) {
	exe := "oracleB.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1571B.go")
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

func genTests() []testB {
	rng := rand.New(rand.NewSource(1571))
	tests := make([]testB, 0, 100)
	for len(tests) < 100 {
		a := rng.Intn(98) + 1
		c := a + 2 + rng.Intn(100-(a+2)+1)
		b := a + 1 + rng.Intn(c-a-1)
		va := rng.Intn(a) + 1
		vc := va + rng.Intn(c-va+1)
		n := c + rng.Intn(100-c+1)
		tests = append(tests, testB{n, a, va, c, vc, b})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		input := fmt.Sprintf("1\n%d\n%d %d\n%d %d\n%d\n", tc.n, tc.a, tc.va, tc.c, tc.vc, tc.b)
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
