package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseD struct {
	n int
	s string
}

func solveD(n int, s string) []int64 {
	improvements := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		if s[i] == 'L' {
			total += int64(i)
			improvements[i] = int64(n - 1 - 2*i)
		} else {
			total += int64(n - 1 - i)
			improvements[i] = int64(2*i - n + 1)
		}
	}
	sort.Slice(improvements, func(i, j int) bool { return improvements[i] > improvements[j] })
	prefix := int64(0)
	ans := make([]int64, n)
	for i := 0; i < n; i++ {
		if improvements[i] > 0 {
			prefix += improvements[i]
		}
		ans[i] = total + prefix
	}
	return ans
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

func generateTests() []testCaseD {
	rand.Seed(45)
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				b[j] = 'L'
			} else {
				b[j] = 'R'
			}
		}
		tests[i] = testCaseD{n, string(b)}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
		expectedSlice := solveD(tc.n, tc.s)
		expected := make([]string, len(expectedSlice))
		for j, v := range expectedSlice {
			expected[j] = fmt.Sprint(v)
		}
		expectedStr := strings.Join(expected, " ")
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(output) != expectedStr {
			fmt.Printf("test %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expectedStr, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
