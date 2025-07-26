package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	x int64
	y int64
}

func expectedA(x, y int64) string {
	if y == 0 {
		return "No"
	}
	if y == 1 {
		if x == 0 {
			return "Yes"
		}
		return "No"
	}
	if x >= y-1 && (x-y+1)%2 == 0 {
		return "Yes"
	}
	return "No"
}

func genTestsA() []testCaseA {
	rand.Seed(1)
	tests := make([]testCaseA, 0, 100)
	// some fixed edge cases
	fixed := []testCaseA{{0, 0}, {0, 1}, {1, 0}, {1, 1}, {2, 1}, {3, 2}, {5, 3}, {10, 1}}
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		x := rand.Int63n(1_000_000_000)
		y := rand.Int63n(1_000_000_000)
		tests = append(tests, testCaseA{x: x, y: y})
	}
	return tests
}

func runCase(bin string, tc testCaseA) error {
	input := fmt.Sprintf("%d %d\n", tc.x, tc.y)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedA(tc.x, tc.y)
	if got != expect {
		return fmt.Errorf("for input %d %d expected %q got %q", tc.x, tc.y, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
