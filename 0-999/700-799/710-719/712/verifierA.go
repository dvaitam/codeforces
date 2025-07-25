package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func solveA(arr []int64) []int64 {
	n := len(arr)
	b := make([]int64, n)
	for i := 0; i < n-1; i++ {
		b[i] = arr[i] + arr[i+1]
	}
	b[n-1] = arr[n-1]
	return b
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		n := rng.Intn(19) + 2 // at least 2
		arr := make([]int64, n)
		var inBuilder strings.Builder
		inBuilder.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			arr[i] = rng.Int63n(2001) - 1000
			if i > 0 {
				inBuilder.WriteByte(' ')
			}
			inBuilder.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		inBuilder.WriteByte('\n')
		res := solveA(arr)
		var outBuilder strings.Builder
		for i := 0; i < n; i++ {
			if i > 0 {
				outBuilder.WriteByte(' ')
			}
			outBuilder.WriteString(fmt.Sprintf("%d", res[i]))
		}
		outBuilder.WriteByte('\n')
		tests[t] = testCase{in: inBuilder.String(), out: outBuilder.String()}
	}
	return tests
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(tc.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
