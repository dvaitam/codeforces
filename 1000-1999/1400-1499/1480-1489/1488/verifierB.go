package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	exe := "oracleB"
	cmd := exec.Command("go", "build", "-o", exe, "1488B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomRBS(n int, rng *rand.Rand) string {
	open := 0
	var sb strings.Builder
	for i := 0; i < n; i++ {
		remaining := n - i
		if open == 0 {
			sb.WriteByte('(')
			open++
			continue
		}
		if open == remaining {
			sb.WriteByte(')')
			open--
			continue
		}
		if open < n/2 && rng.Intn(2) == 0 {
			sb.WriteByte('(')
			open++
		} else {
			sb.WriteByte(')')
			open--
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(10)*2 + 2
		k := rng.Intn(n) + 1
		s := randomRBS(n, rng)
		fmt.Fprintf(&sb, "%d %d\n%s\n", n, k, s)
	}
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Printf("oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d mismatch\nexpected:\n%s\n got:\n%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
