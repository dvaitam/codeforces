package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// solve computes expected output for input containing two integers a and b.
func solve(input string) string {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) < 2 {
		return ""
	}
	var a, b int
	fmt.Sscan(parts[0], &a)
	fmt.Sscan(parts[1], &b)
	win1, draw, win2 := 0, 0, 0
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	for x := 1; x <= 6; x++ {
		da := abs(a - x)
		db := abs(b - x)
		if da < db {
			win1++
		} else if da == db {
			draw++
		} else {
			win2++
		}
	}
	return fmt.Sprintf("%d %d %d", win1, draw, win2)
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type test struct {
	input    string
	expected string
}

// generateTests creates at least 100 test cases
func generateTests() []test {
	rand.Seed(42)
	var tests []test
	for a := 1; a <= 6; a++ {
		for b := 1; b <= 6; b++ {
			inp := fmt.Sprintf("%d %d\n", a, b)
			tests = append(tests, test{inp, solve(inp)})
		}
	}
	for len(tests) < 100 {
		a := rand.Intn(6) + 1
		b := rand.Intn(6) + 1
		inp := fmt.Sprintf("%d %d\n", a, b)
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
