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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
	exp, _ := strconv.Atoi(expect)
	val, err := strconv.Atoi(actual)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([]int{2, 2, 2, 2}),
		makeCase([]int{5, 1, 4, 5, 1, 2}),
		makeCase([]int{5, 4, 2, 1, 2, 3}),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(n) + 1
		}
		tests = append(tests, makeCase(arr))
	}
	return tests
}

func makeCase(arr []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	expect := reference(arr)
	return testCase{
		input:  sb.String(),
		expect: strconv.Itoa(expect),
	}
}

func reference(arr []int) int {
	seen := make(map[int]struct{})
	for _, v := range arr {
		seen[v] = struct{}{}
	}
	return len(seen)
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
