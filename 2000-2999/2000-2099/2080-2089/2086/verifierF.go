package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	refSourceF = "2086F.go"
	refBinaryF = "refF.bin"
)

type testCase struct {
	n int
	s string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if err := validate(tc, refOut); err != nil {
			fmt.Printf("reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, refOut)
			return
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if err := validate(tc, candOut); err != nil {
			fmt.Printf("candidate failed on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, candOut)
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryF, refSourceF)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryF), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func validate(tc testCase, output string) error {
	tokens := strings.Fields(output)
	expectedTokens := tc.n * 2
	if len(tokens) != expectedTokens {
		return fmt.Errorf("expected %d integers (n pairs), got %d", expectedTokens, len(tokens))
	}

	t := make([]byte, 0, tc.n)
	for step := 0; step < tc.n; step++ {
		t = append(t, tc.s[step])
		l, err := strconv.Atoi(tokens[2*step])
		if err != nil {
			return fmt.Errorf("token %d is not an integer: %v", 2*step+1, err)
		}
		r, err := strconv.Atoi(tokens[2*step+1])
		if err != nil {
			return fmt.Errorf("token %d is not an integer: %v", 2*step+2, err)
		}

		if l == 0 && r == 0 {
			continue
		}
		if l == 0 || r == 0 {
			return fmt.Errorf("invalid zero index at step %d: (%d, %d)", step+1, l, r)
		}
		if l < 1 || l > step+1 || r < 1 || r > step+1 {
			return fmt.Errorf("swap indices out of range at step %d: (%d, %d), current length %d", step+1, l, r, step+1)
		}
		t[l-1], t[r-1] = t[r-1], t[l-1]
	}

	if !isPalindrome(t) {
		return fmt.Errorf("final string %q is not a palindrome", string(t))
	}
	return nil
}

func isPalindrome(b []byte) bool {
	for l, r := 0, len(b)-1; l < r; l, r = l+1, r-1 {
		if b[l] != b[r] {
			return false
		}
	}
	return true
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	sb.WriteString(tc.s)
	if !strings.HasSuffix(tc.s, "\n") {
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2086))
	var tests []testCase

	addCase := func(s string) {
		tests = append(tests, testCase{n: len(s), s: s})
	}

	addCase("a")
	addCase("b")
	addCase(strings.Repeat("a", 99))
	addCase(strings.Repeat("b", 99))
	addCase("ababa")
	addCase("babab")

	for len(tests) < 200 {
		n := rnd.Intn(50)*2 + 1 // odd length between 1 and 99
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rnd.Intn(2) == 0 {
				sb.WriteByte('a')
			} else {
				sb.WriteByte('b')
			}
		}
		addCase(sb.String())
	}
	return tests
}
