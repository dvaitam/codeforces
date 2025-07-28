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
	oracle := filepath.Join(dir, "oracleE3")
	cmd := exec.Command("go", "build", "-o", oracle, "1808E3.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

var primes = []int64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73}

func generateCase(rng *rand.Rand) string {
	n := rng.Int63n(50) + 1
	k := rng.Int63n(8) + 2
	m := primes[rng.Intn(len(primes))]
	return fmt.Sprintf("%d %d %d\n", n, k, m)
}

func runCase(bin, oracle, input string) error {
	run := func(exe string) (string, error) {
		var cmd *exec.Cmd
		if strings.HasSuffix(exe, ".go") {
			cmd = exec.Command("go", "run", exe)
		} else {
			cmd = exec.Command(exe)
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
	exp, err := run(oracle)
	if err != nil {
		return fmt.Errorf("oracle: %v", err)
	}
	got, err := run(bin)
	if err != nil {
		return err
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE3.go /path/to/binary")
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
		in := generateCase(rng)
		if err := runCase(bin, oracle, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
