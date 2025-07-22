package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "394E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
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

func generate() []string {
	const T = 100
	rand.Seed(5)
	tests := make([]string, T)
	for t := 0; t < T; t++ {
		n := rand.Intn(4) + 2
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			x := rand.Intn(21) - 5
			y := rand.Intn(21) - 5
			fmt.Fprintf(&sb, "%d %d\n", x, y)
		}
		sb.WriteString("4\n")
		sb.WriteString("0 0\n0 10\n10 10\n10 0\n")
		tests[t] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := generate()
	for idx, input := range tests {
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
