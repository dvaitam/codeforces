package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	a string
	b string
}

func solve(a, b string) []string {
	if a == b {
		return []string{"YES", a}
	}
	if a[0] == b[0] {
		return []string{"YES", fmt.Sprintf("%c*", a[0])}
	}
	if a[len(a)-1] == b[len(b)-1] {
		return []string{"YES", fmt.Sprintf("*%c", a[len(a)-1])}
	}
	for i := 0; i+1 < len(a); i++ {
		sub := a[i : i+2]
		if strings.Contains(b, sub) {
			return []string{"YES", "*" + sub + "*"}
		}
	}
	return []string{"NO"}
}

func genTests() []TestCase {
	letters := "abcde"
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		la := i%5 + 1
		lb := (i/5)%5 + 1
		var sb strings.Builder
		for j := 0; j < la; j++ {
			sb.WriteByte(letters[(i+j)%len(letters)])
		}
		a := sb.String()
		sb.Reset()
		for j := 0; j < lb; j++ {
			sb.WriteByte(letters[(i+2*j)%len(letters)])
		}
		b := sb.String()
		tests = append(tests, TestCase{a, b})
	}
	return tests
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.a)
		fmt.Fprintln(&input, tc.b)
	}

	expectedLines := []string{}
	for _, tc := range tests {
		expectedLines = append(expectedLines, solve(tc.a, tc.b)...)
	}

	out, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		fmt.Print(out)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(expectedLines) {
		fmt.Printf("wrong number of lines: got %d want %d\n", len(lines), len(expectedLines))
		os.Exit(1)
	}
	for i, got := range lines {
		got = strings.TrimSpace(got)
		if got != expectedLines[i] {
			tcIndex := 0
			count := 0
			for j := 0; j < len(expectedLines); j++ {
				if count == i {
					tcIndex = j
					break
				}
				if expectedLines[j] == "YES" || expectedLines[j] == "NO" {
					count++
				}
			}
			fmt.Printf("mismatch on line %d: expected %s got %s\n", i+1, expectedLines[i], got)
			fmt.Printf("corresponding test #%d: %s %s\n", tcIndex+1, tests[tcIndex].a, tests[tcIndex].b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
