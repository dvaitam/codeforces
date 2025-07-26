package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "835C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runOracle(oracle, input string) (string, error) {
	cmd := exec.Command(oracle)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	c := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, q, c)
	for i := 0; i < n; i++ {
		x := rng.Intn(100) + 1
		y := rng.Intn(100) + 1
		s := rng.Intn(c + 1)
		fmt.Fprintf(&sb, "%d %d %d\n", x, y, s)
	}
	for i := 0; i < q; i++ {
		t := rng.Intn(101)
		x1 := rng.Intn(100) + 1
		y1 := rng.Intn(100) + 1
		x2 := x1 + rng.Intn(101-x1)
		y2 := y1 + rng.Intn(101-y1)
		if x2 == 0 {
			x2 = 1
		}
		if y2 == 0 {
			y2 = 1
		}
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", t, x1, y1, x2, y2)
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	expected, err := runOracle(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
