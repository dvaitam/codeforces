package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
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

func gcd(a, b *big.Int) *big.Int {
	var g big.Int
	return g.GCD(nil, nil, a, b)
}

func checkSolution(n int, tokens []string) error {
	if n == 2 {
		if len(tokens) != 1 || strings.TrimSpace(tokens[0]) != "-1" {
			return fmt.Errorf("expected -1 for n=2")
		}
		return nil
	}
	if len(tokens) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	nums := make([]*big.Int, n)
	seen := make(map[string]struct{})
	for i, t := range tokens {
		if len(t) > 100 {
			return fmt.Errorf("number %d too long", i+1)
		}
		if strings.HasPrefix(t, "0") && t != "0" {
			return fmt.Errorf("number %d has leading zeros", i+1)
		}
		z, ok := new(big.Int).SetString(t, 10)
		if !ok || z.Sign() <= 0 {
			return fmt.Errorf("invalid integer %q", t)
		}
		nums[i] = z
		if _, ok := seen[z.String()]; ok {
			return fmt.Errorf("numbers not distinct")
		}
		seen[z.String()] = struct{}{}
	}
	totalGCD := new(big.Int).Set(nums[0])
	for i := 1; i < n; i++ {
		totalGCD = gcd(totalGCD, nums[i])
	}
	if totalGCD.Cmp(big.NewInt(1)) != 0 {
		return fmt.Errorf("gcd of all numbers not 1")
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if gcd(new(big.Int).Set(nums[i]), nums[j]).Cmp(big.NewInt(1)) == 0 {
				return fmt.Errorf("gcd of pair %d,%d is 1", i+1, j+1)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		n := 0
		fmt.Sscan(line, &n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		outTokens := strings.Fields(got)
		if err := checkSolution(n, outTokens); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
