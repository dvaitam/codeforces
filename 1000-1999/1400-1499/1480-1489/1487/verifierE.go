package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "1487E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase(rng *rand.Rand) string {
	n1 := rng.Intn(5) + 1
	n2 := rng.Intn(5) + 1
	n3 := rng.Intn(5) + 1
	n4 := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n1, n2, n3, n4))
	for i := 0; i < n1; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n2; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n3; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n4; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	m1 := rng.Intn(n1*n2 + 1)
	sb.WriteString(fmt.Sprintf("%d\n", m1))
	for i := 0; i < m1; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", rng.Intn(n1)+1, rng.Intn(n2)+1))
	}
	m2 := rng.Intn(n2*n3 + 1)
	sb.WriteString(fmt.Sprintf("%d\n", m2))
	for i := 0; i < m2; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", rng.Intn(n2)+1, rng.Intn(n3)+1))
	}
	m3 := rng.Intn(n3*n4 + 1)
	sb.WriteString(fmt.Sprintf("%d\n", m3))
	for i := 0; i < m3; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", rng.Intn(n3)+1, rng.Intn(n4)+1))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(1))
    for t := 0; t < 100; t++ {
        tc := generateCase(rng)
        input := tc
        exp, err := run(oracle, input)
        if err != nil {
            fmt.Printf("oracle runtime error: %v\n", err)
            os.Exit(1)
        }
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\ninput:\n%sexpected %s got %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
