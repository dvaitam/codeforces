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
	a, b, c int
}

func generateCases() []testCase {
	rand.Seed(1)
	cases := make([]testCase, 100)
	for i := range cases {
		cases[i] = testCase{rand.Intn(21), rand.Intn(21), rand.Intn(21)}
	}
	return cases
}

func expected(t testCase) string {
	if t.a+t.b == t.c || t.a+t.c == t.b || t.b+t.c == t.a {
		return "YES"
	}
	return "NO"
}

func buildIO(cases []testCase) (string, string) {
	var inBuilder strings.Builder
	var outBuilder strings.Builder
	fmt.Fprintf(&inBuilder, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&inBuilder, "%d %d %d\n", tc.a, tc.b, tc.c)
		fmt.Fprintln(&outBuilder, expected(tc))
	}
	return inBuilder.String(), outBuilder.String()
}

func run(binary, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(binary, ".go") {
		cmd = exec.Command("go", "run", binary)
	} else {
		cmd = exec.Command(binary)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	cases := generateCases()
	input, expectedOutput := buildIO(cases)
	actualOutput, err := run(binary, input)
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		os.Exit(1)
	}
	if normalize(actualOutput) != normalize(expectedOutput) {
		fmt.Println("Wrong answer")
		fmt.Println("Expected:")
		fmt.Println(expectedOutput)
		fmt.Println("Got:")
		fmt.Println(actualOutput)
		os.Exit(1)
	}
	fmt.Println("All test cases passed!")
}
