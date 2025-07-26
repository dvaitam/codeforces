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
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1271F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		a1 := rng.Intn(10) + 1
		b1 := rng.Intn(10) + 1
		c1 := rng.Intn(10) + 1
		a2 := rng.Intn(10) + 1
		b2 := rng.Intn(10) + 1
		c2 := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", a1, b1, c1)
		fmt.Fprintf(&sb, "%d %d %d\n", a2, b2, c2)
		for j := 0; j < 7; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			d := rng.Intn(5)
			fmt.Fprintf(&sb, "%d", d)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func run(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
