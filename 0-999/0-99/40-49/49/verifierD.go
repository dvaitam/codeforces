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
	n int
	s string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n%s\n", t.n, t.s)
		expect := solveD(t.n, t.s)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveD(n int, s string) string {
	ans := 0
	for i := 0; i < n; {
		j := i + 1
		for j < n && s[j] == s[i] {
			j++
		}
		ans += (j - i) / 2
		i = j
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateTests() []testCase {
	rand.Seed(4)
	tests := make([]testCase, 0, 100)
	fixed := []string{"0", "1", "01", "10", "0000", "1111", "0101", "1010"}
	for _, f := range fixed {
		tests = append(tests, testCase{len(f), f})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests = append(tests, testCase{n, sb.String()})
	}
	return tests
}
