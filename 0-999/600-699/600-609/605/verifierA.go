package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, filepath.Join(dir, "605A.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	perm := rng.Perm(n)
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i, v := range perm {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected: %s\nactual: %s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
