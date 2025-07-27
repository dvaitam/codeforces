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
	s string
}

func genTestsA() []testCaseA {
	rand.Seed(42)
	tests := make([]testCaseA, 100)
	for i := range tests {
		n := rand.Intn(100) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rand.Intn(3) == 0 {
				sb.WriteByte(')')
			} else {
				sb.WriteByte(byte('a' + rand.Intn(26)))
			}
		}
		tests[i] = testCaseA{sb.String()}
	}
	// add edge cases
	tests = append(tests, testCaseA{")"})
	tests = append(tests, testCaseA{"abc"})
	return tests
}

func solveA(tc testCaseA) string {
	n := len(tc.s)
	cnt := 0
	for i := n - 1; i >= 0 && tc.s[i] == ')'; i-- {
		cnt++
	}
	if cnt > n-cnt {
		return "Yes"
	}
	return "No"
}

func runBinary(bin string, input string) (string, error) {
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
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n", len(tc.s), tc.s)
		exp := solveA(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
