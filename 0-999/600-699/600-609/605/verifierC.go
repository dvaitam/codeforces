package main

import (
	"bytes"
	"fmt"
	"math"
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
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, filepath.Join(dir, "605C.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	p := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, p, q)
	for i := 0; i < n; i++ {
		a := rng.Intn(20) + 1
		d := rng.Intn(20) + 1
		fmt.Fprintf(&b, "%d %d\n", a, d)
	}
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

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscan(strings.TrimSpace(s), &f)
	return f, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		expect, _ := parseFloat(expectStr)
		got, _ := parseFloat(gotStr)
		if math.Abs(expect-got) > 1e-6*math.Max(1, math.Abs(expect)) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected: %s\nactual: %s\n", i+1, input, expectStr, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
