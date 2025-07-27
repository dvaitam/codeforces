package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n int64
}

func genTestsB() []testCaseB {
	rand.Seed(43)
	tests := make([]testCaseB, 100)
	for i := range tests {
		var n int64
		if rand.Intn(2) == 0 {
			n = int64(rand.Intn(1000) + 1)
		} else {
			n = rand.Int63n(1e12) + 1
		}
		tests[i] = testCaseB{n}
	}
	tests = append(tests, testCaseB{1})
	tests = append(tests, testCaseB{282})
	return tests
}

func isFair(x int64) bool {
	t := x
	for t > 0 {
		d := t % 10
		if d != 0 && x%int64(d) != 0 {
			return false
		}
		t /= 10
	}
	return true
}

func solveB(tc testCaseB) int64 {
	x := tc.n
	for {
		if isFair(x) {
			return x
		}
		x++
	}
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GOMAXPROCS=1")
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n", tc.n)
		exp := solveB(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != fmt.Sprintf("%d", exp) {
			fmt.Printf("test %d failed: expected %d got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
