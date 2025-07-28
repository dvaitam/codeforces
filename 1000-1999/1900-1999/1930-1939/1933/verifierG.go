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

const numTestsG = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "binG")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleG")
	cmd := exec.Command("go", "build", "-o", tmp, "1933G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func randShape(rng *rand.Rand) string {
	if rng.Intn(2) == 0 {
		return "circle"
	}
	return "square"
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 5
	m := rng.Intn(3) + 5
	q := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, m, q)
	for i := 0; i < q; i++ {
		r := rng.Intn(n) + 1
		c := rng.Intn(m) + 1
		fmt.Fprintf(&sb, "%d %d %s\n", r, c, randShape(rng))
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	expected, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		return
	}
	path := os.Args[len(os.Args)-1]
	bin, cleanup, err := prepareBinary(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer cleanOracle()
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < numTestsG; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
