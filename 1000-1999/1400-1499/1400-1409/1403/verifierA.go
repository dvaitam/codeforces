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
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "1403A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	N := rng.Intn(5) + 2
	D := rng.Intn(5) + 1
	U := rng.Intn(4) + 1
	Q := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", N, D, U, Q))
	for i := 0; i < N; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(20)))
	}
	sb.WriteByte('\n')
	for i := 0; i < U; i++ {
		a := rng.Intn(N)
		b := rng.Intn(N - 1)
		if b >= a {
			b++
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	for i := 0; i < Q; i++ {
		x := rng.Intn(N)
		y := rng.Intn(N)
		v := rng.Intn(U + 1)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, v))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
