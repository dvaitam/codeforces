package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// buildOracle compiles 163B.go into a temporary binary.
func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "163B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

// run executes a binary or go source file with given input.
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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	h := rng.Intn(10) + 1
	weights := make([]int, n)
	speeds := make([]int, n)
	for i := 0; i < n; i++ {
		weights[i] = rng.Intn(50) + 1
	}
	for i := 0; i < n; i++ {
		speeds[i] = rng.Intn(50) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, h))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d", weights[i]))
		if i+1 < n {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("\n")
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d", speeds[i]))
		if i+1 < n {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("\n")
	return sb.String()
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
		input := generateCase(rng)
		expected, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
