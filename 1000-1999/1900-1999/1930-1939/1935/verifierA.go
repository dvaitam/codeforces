package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	n int
	s string
}

func reverseString(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func expected(tc TestCase) string {
	rev := reverseString(tc.s)
	ans := tc.s
	if tmp := tc.s + rev; tmp < ans {
		ans = tmp
	}
	if tmp := rev + tc.s; tmp < ans {
		ans = tmp
	}
	return ans
}

func genTests() []TestCase {
	rand.Seed(time.Now().UnixNano())
	const T = 100
	tests := make([]TestCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(10) + 2
		if n%2 == 1 {
			n++
		}
		l := rand.Intn(10) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rand.Intn(26))
		}
		tests[i] = TestCase{n: n, s: string(b)}
	}
	return tests
}

func buildInput(tests []TestCase) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&buf, tc.n)
		fmt.Fprintln(&buf, tc.s)
	}
	return buf.Bytes()
}

func runBinary(path string, input []byte) ([]byte, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	return cmd.CombinedOutput()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := genTests()
	input := buildInput(tests)

	expectedOutputs := make([]string, len(tests))
	for i, tc := range tests {
		expectedOutputs[i] = expected(tc)
	}

	out, err := runBinary(binary, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}

	got := strings.Fields(strings.TrimSpace(string(out)))
	if len(got) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(expectedOutputs), len(got))
		os.Exit(1)
	}
	for i := range expectedOutputs {
		if got[i] != expectedOutputs[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\ninput:\n%d\n%s\nexpected: %s\nactual: %s\n", i+1, tests[i].n, tests[i].s, expectedOutputs[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
