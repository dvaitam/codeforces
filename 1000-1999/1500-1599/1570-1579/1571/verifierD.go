package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testD struct {
	n     int
	m     int
	pairs [][2]int
}

func buildOracle() (string, error) {
	exe := "oracleD.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1571D.go")
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

func genTests() []testD {
	rng := rand.New(rand.NewSource(1571))
	tests := make([]testD, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(5) + 2
		m := rng.Intn(6) + 1
		pairs := make([][2]int, m)
		for i := 0; i < m; i++ {
			f := rng.Intn(n) + 1
			l := rng.Intn(n-1) + 1
			if l >= f {
				l++
			}
			pairs[i] = [2]int{f, l}
		}
		tests = append(tests, testD{n, m, pairs})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, p := range tc.pairs {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		input := sb.String()
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
