package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(tc.expect, got); err != nil {
			fmt.Printf("test %d failed: %v\nexpected: %q\nactual: %q\n", i+1, err, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	if strings.TrimSpace(actual) != expect {
		return fmt.Errorf("expected no output but received non-empty output")
	}
	return nil
}

func genTests() []testCase {
	return []testCase{
		{input: "", expect: solveRef("")},
	}
}

func solveRef(_ string) string {
	return ""
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
