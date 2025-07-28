package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
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
	cmd := exec.Command("go", "build", "-o", oracle, "1487F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randBig(rng *rand.Rand) string {
	digits := rng.Intn(18) + 1
	var sb strings.Builder
	for i := 0; i < digits; i++ {
		d := rng.Intn(10)
		if i == 0 && d == 0 {
			d = 1
		}
		sb.WriteByte(byte('0' + d))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
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
		s := randBig(rng)
		input := s + "\n"
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
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed:\ninput:%sexpected %s got %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
		// ensure output is numeric
		if _, ok := new(big.Int).SetString(strings.TrimSpace(got), 10); !ok {
			fmt.Printf("case %d: output not integer: %s\n", t+1, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
