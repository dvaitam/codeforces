package main

import (
	"bytes"
	"fmt"
	"math"
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
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "185B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) string {
	s := rng.Intn(1000) + 1
	a := rng.Intn(1001)
	b := rng.Intn(1001)
	c := rng.Intn(1001)
	return fmt.Sprintf("%d\n%d %d %d\n", s, a, b, c)
}

func runCase(bin, oracle, input string) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())

	expectedFloats, err := parseFloats(expected)
	if err != nil {
		return fmt.Errorf("parse expected output: %v", err)
	}
	gotFloats, err := parseFloats(got)
	if err != nil {
		return fmt.Errorf("parse contestant output: %v", err)
	}
	if !equalWithTolerance(expectedFloats, gotFloats, 1e-9) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func parseFloats(s string) ([]float64, error) {
	fields := strings.Fields(s)
	res := make([]float64, len(fields))
	for i, f := range fields {
		v, err := parseFloat(f)
		if err != nil {
			return nil, fmt.Errorf("field %d (%q): %w", i, f, err)
		}
		res[i] = v
	}
	return res, nil
}

func parseFloat(s string) (float64, error) {
	var v float64
	if _, err := fmt.Sscan(s, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func equalWithTolerance(expected, got []float64, eps float64) bool {
	if len(expected) != len(got) {
		return false
	}
	for i := range expected {
		if math.Abs(expected[i]-got[i]) > eps {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
