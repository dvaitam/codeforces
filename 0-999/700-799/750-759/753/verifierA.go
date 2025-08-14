package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) int {
	return rng.Intn(1000) + 1
}

func runProg(exe, input string) (string, error) {
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

func maxK(n int) int {
	k := 0
	for (k+1)*(k+2)/2 <= n {
		k++
	}
	return k
}

func checkOutput(n int, out string, kExpected int) error {
	r := strings.NewReader(out)
	var k int
	if _, err := fmt.Fscan(r, &k); err != nil {
		return fmt.Errorf("failed to read k: %v", err)
	}
	var nums []int
	for {
		var x int
		_, err := fmt.Fscan(r, &x)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("invalid number: %v", err)
		}
		nums = append(nums, x)
	}
	if len(nums) != k {
		return fmt.Errorf("expected %d numbers, got %d", k, len(nums))
	}
	if k != kExpected {
		return fmt.Errorf("expected k=%d, got %d", kExpected, k)
	}
	sum := 0
	seen := make(map[int]bool)
	for _, v := range nums {
		if v <= 0 {
			return fmt.Errorf("non-positive number %d", v)
		}
		if seen[v] {
			return fmt.Errorf("number %d repeated", v)
		}
		seen[v] = true
		sum += v
	}
	if sum != n {
		return fmt.Errorf("sum mismatch: expected %d, got %d", n, sum)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := generateCase(rng)
		input := fmt.Sprintf("%d\n", n)
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkOutput(n, got, maxK(n)); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
