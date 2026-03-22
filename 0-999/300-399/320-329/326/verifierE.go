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
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "326E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(bin, oracle, name string, n int64, h int) error {
	input := fmt.Sprintf("%s\n%d %d\n", name, n, h)

	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	var exp float64
	if _, err := fmt.Sscan(strings.TrimSpace(outO.String()), &exp); err != nil {
		return fmt.Errorf("oracle output parse error: %v", err)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var val float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &val); err != nil {
		return fmt.Errorf("unable to parse output: %v", err)
	}

	tol := math.Max(1e-6, abs(exp)*1e-9)
	if abs(val-exp) > tol {
		return fmt.Errorf("expected %.10f got %.10f", exp, val)
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int64, int) {
	name := "Alice"
	if rng.Intn(2) == 0 {
		name = "Bob"
	}
	return name, rng.Int63n(1000), rng.Intn(10)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		name, n, h := generateCase(rng)
		if err := runCase(bin, oracle, name, n, h); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
