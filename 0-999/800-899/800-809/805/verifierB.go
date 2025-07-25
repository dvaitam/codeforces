package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n int
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(1)
	tests := make([]Test, 0, 150)
	for i := 1; i <= 50; i++ {
		tests = append(tests, Test{i})
	}
	for i := 0; i < 100; i++ {
		tests = append(tests, Test{rand.Intn(100) + 1})
	}
	return tests
}

func check(n int, out string) error {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return fmt.Errorf("expected single token output")
	}
	s := fields[0]
	if len(s) != n {
		return fmt.Errorf("expected length %d got %d", n, len(s))
	}
	for i := 0; i < n; i++ {
		c := s[i]
		if c != 'a' && c != 'b' && c != 'c' {
			return fmt.Errorf("invalid character %q", c)
		}
		if c == 'c' {
			return fmt.Errorf("string contains 'c'")
		}
		if i+2 < n && s[i] == s[i+2] {
			return fmt.Errorf("palindrome of length 3 at position %d", i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.n)
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(tc.n, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:%soutput:%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
