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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	val, err := strconv.ParseInt(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output not an integer: %v", err)
	}
	exp, _ := strconv.ParseInt(expect, 10, 64)
	if val != exp {
		return fmt.Errorf("expected %s but got %s", expect, actual)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	var tests []testCase
	tests = append(tests, makeTest([]int{0}))
	tests = append(tests, makeTest([]int{1, 2, 3}))
	tests = append(tests, makeTest([]int{3979, 3979, 3979}))
	for i := 0; i < 200; i++ {
		n := rand.Intn(100) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(3980)
		}
		tests = append(tests, makeTest(arr))
	}
	return tests
}

func makeTest(arr []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", sum),
	}
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
