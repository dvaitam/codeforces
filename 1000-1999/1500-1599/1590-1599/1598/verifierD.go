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

type testCase struct {
	input    string
	expected int64
}

func expectedD(a, b []int) int64 {
	n := len(a)
	cntA := make(map[int]int)
	cntB := make(map[int]int)
	for i := 0; i < n; i++ {
		cntA[a[i]]++
		cntB[b[i]]++
	}
	total := int64(n) * (int64(n) - 1) * (int64(n) - 2) / 6
	var bad int64
	for i := 0; i < n; i++ {
		bad += int64(cntA[a[i]]-1) * int64(cntB[b[i]]-1)
	}
	return total - bad
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 3
		arrA := make([]int, n)
		arrB := make([]int, n)
		
		seen := make(map[string]bool)
		for j := 0; j < n; j++ {
			for {
				a := rng.Intn(n) + 1
				b := rng.Intn(n) + 1
				key := fmt.Sprintf("%d,%d", a, b)
				if !seen[key] {
					seen[key] = true
					arrA[j] = a
					arrB[j] = b
					break
				}
			}
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", arrA[j], arrB[j]))
		}
		exp := expectedD(arrA, arrB)
		cases[i] = testCase{input: sb.String(), expected: exp}
	}
	return cases
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Printf("usage: go run verifierD.go /path/to/binary. Got args: %v\n", os.Args)
		os.Exit(1)
	}
	bin := args[0]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(out, &got)
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
