package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseB struct {
	n  int
	s1 string
	s2 string
}

func solveB(n int, s1, s2 string) string {
	for i := 0; i < n; i++ {
		a := s1[i] == 'R'
		b := s2[i] == 'R'
		if a != b {
			return "NO"
		}
	}
	return "YES"
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseB {
	rand.Seed(43)
	tests := make([]testCaseB, 100)
	letters := []byte{'R', 'G', 'B'}
	for i := range tests {
		n := rand.Intn(20) + 1
		b1 := make([]byte, n)
		b2 := make([]byte, n)
		for j := 0; j < n; j++ {
			b1[j] = letters[rand.Intn(3)]
			b2[j] = letters[rand.Intn(3)]
		}
		tests[i] = testCaseB{n, string(b1), string(b2)}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", tc.n, tc.s1, tc.s2)
		expected := solveB(tc.n, tc.s1, tc.s2)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("test %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expected, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
