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
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "695D.go")
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

func randField(max int, rng *rand.Rand) int {
	if rng.Float64() < 0.3 {
		return -1
	}
	return rng.Intn(max + 1)
}

func generateCase(rng *rand.Rand) string {
	s := randField(59, rng)
	m := randField(59, rng)
	h := randField(23, rng)
	dow := randField(7, rng)
	if dow != -1 {
		dow = rng.Intn(7) + 1
	}
	dom := randField(31, rng)
	if dom != -1 {
		dom = rng.Intn(31) + 1
	}
	mo := randField(12, rng)
	if mo != -1 {
		mo = rng.Intn(12) + 1
	}
	n := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d %d %d\n", s, m, h, dow, dom, mo)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		t := rng.Int63n(1_000_000)
		fmt.Fprintf(&sb, "%d\n", t)
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
