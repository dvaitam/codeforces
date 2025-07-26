package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseC struct {
	n uint64
	k uint64
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedC(n, k uint64) string {
	l := uint64(1)
	limit := n + 1
	for i := uint64(2); i <= k && l <= limit; i++ {
		g := gcd(l, i)
		if l/g > limit/i {
			l = limit + 1
			break
		}
		l = l / g * i
	}
	if l <= limit && limit%l == 0 {
		return "Yes"
	}
	return "No"
}

func genTestsC() []testCaseC {
	rand.Seed(3)
	tests := make([]testCaseC, 0, 100)
	fixed := []testCaseC{{1, 1}, {10, 1}, {10, 2}, {10, 5}, {100, 10}}
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		n := uint64(rand.Int63n(1_000_000_000_000)) + 1
		k := uint64(rand.Intn(5000) + 1)
		tests = append(tests, testCaseC{n: n, k: k})
	}
	return tests
}

func runCase(bin string, tc testCaseC) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
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
	expect := expectedC(tc.n, tc.k)
	if got != expect {
		return fmt.Errorf("n=%d k=%d expected %s got %s", tc.n, tc.k, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
