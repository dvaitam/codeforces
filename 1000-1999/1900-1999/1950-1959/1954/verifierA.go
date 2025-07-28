package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveCase(n, m, k int) string {
	if k >= n {
		return "NO"
	}
	maxPerColor := (n + m - 1) / m
	if maxPerColor < n-k {
		return "YES"
	}
	return "NO"
}

func generateTests() []string {
	rand.Seed(1)
	tests := make([]string, 100)
	for i := range tests {
		n := rand.Intn(50) + 1
		m := rand.Intn(n) + 1
		k := rand.Intn(n) + 1
		tests[i] = fmt.Sprintf("%d %d %d", n, m, k)
	}
	return tests
}

func expectedOutputs(cases []string) []string {
	out := make([]string, len(cases))
	for i, c := range cases {
		var n, m, k int
		fmt.Sscanf(c, "%d %d %d", &n, &m, &k)
		out[i] = solveCase(n, m, k)
	}
	return out
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	b, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return string(b), fmt.Errorf("time limit exceeded")
	}
	return string(b), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := generateTests()
	tInput := fmt.Sprintf("%d\n%s\n", len(cases), strings.Join(cases, "\n"))
	exp := expectedOutputs(cases)

	out, err := runBinary(bin, tInput)
	if err != nil {
		fmt.Println("Error running binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	var got []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			got = append(got, strings.ToUpper(line))
		}
	}
	if len(got) != len(exp) {
		fmt.Printf("Mismatch number of lines: expected %d got %d\n", len(exp), len(got))
		os.Exit(1)
	}
	for i := range exp {
		if got[i] != exp[i] {
			fmt.Printf("Test %d failed: input %s expected %s got %s\n", i+1, cases[i], exp[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
