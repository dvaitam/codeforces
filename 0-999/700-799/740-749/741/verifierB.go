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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "741B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	maxM := n * (n - 1) / 2
	m := rng.Intn(maxM + 1)
	W := rng.Intn(20) + 1
	w := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		w[i] = rng.Intn(10) + 1
		b[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, W))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(w[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		in := generateCase(rng)
		exp, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
