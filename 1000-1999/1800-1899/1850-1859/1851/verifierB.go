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
	oracle := filepath.Join(os.TempDir(), "oracleB1851.bin")
	cmd := exec.Command("go", "build", "-o", oracle, "1851B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generate() string {
	const T = 100
	rng := rand.New(rand.NewSource(2))
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", T)
	for i := 0; i < T; i++ {
		n := rng.Intn(50) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(1000)+1)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	input := generate()
	exp, err := run(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	got, err := run(cand, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(exp) != strings.TrimSpace(got) {
		fmt.Println("wrong answer")
		fmt.Println("input:\n" + input)
		fmt.Println("expected:\n" + exp)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
