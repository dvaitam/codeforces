package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	A int
	B int
	N int
}

func digits(x int) int {
	if x == 0 {
		return 1
	}
	d := 0
	for x > 0 {
		d++
		x /= 10
	}
	return d
}

func solve(tc TestCase) int {
	A := tc.A
	B := tc.B
	N := tc.N
	pow10 := make([]int, 11)
	pow10[0] = 1
	for i := 1; i < len(pow10); i++ {
		pow10[i] = pow10[i-1] * 10
	}
	maxK := 1
	for k := 1; k < len(pow10); k++ {
		if pow10[k] <= N {
			maxK = k
		}
	}
	count := 0
	for a := 1; a < A; a++ {
		for b := 1; b < B; b++ {
			m := digits(b)
			den := a*pow10[m] - b
			if den <= 0 {
				continue
			}
			for k := 1; k <= maxK; k++ {
				num := a * b * (pow10[k] - 1)
				if num%den != 0 {
					continue
				}
				n := num / den
				if n >= 1 && n < N && digits(n) == k {
					count++
				}
			}
		}
	}
	return count
}

func genTests() []TestCase {
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		A := 3 + i%5
		B := 3 + (i/5)%5
		N := 5 + i*3
		tests = append(tests, TestCase{A, B, N})
	}
	return tests
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	for i, tc := range tests {
		expected := fmt.Sprintf("%d", solve(tc))
		input := fmt.Sprintf("%d %d %d\n", tc.A, tc.B, tc.N)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error running binary:", err)
			fmt.Print(out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != expected {
			fmt.Printf("test %d failed expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
