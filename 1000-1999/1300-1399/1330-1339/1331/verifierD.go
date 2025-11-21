package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	if actual == "" {
		return fmt.Errorf("empty output")
	}
	got, err := strconv.ParseUint(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not a valid unsigned integer: %v", err)
	}
	exp, _ := strconv.ParseUint(expect, 10, 64)
	if got != exp {
		return fmt.Errorf("expected %s but got %s", expect, actual)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest("A000000"),
		makeTest("A123456"),
		makeTest("A999999"),
	}
	for i := 0; i < 200; i++ {
		tests = append(tests, makeTest(randomHex()))
	}
	return tests
}

func makeTest(hex string) testCase {
	value, _ := strconv.ParseUint(hex, 16, 64)
	return testCase{
		input:  hex + "\n",
		expect: fmt.Sprintf("%d", value),
	}
}

func randomHex() string {
	var sb strings.Builder
	sb.WriteByte('A')
	for i := 0; i < 6; i++ {
		sb.WriteByte(byte('0' + rand.Intn(10)))
	}
	return sb.String()
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
