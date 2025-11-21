package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input string
	array []int
}

type operation struct {
	i, j int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
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
		if err := check(tc, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(tc testCase, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("unable to read number of operations: %v", err)
	}
	if k < 0 || k > 50 {
		return fmt.Errorf("invalid number of operations %d", k)
	}
	ops := make([]operation, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &ops[i].i, &ops[i].j); err != nil {
			return fmt.Errorf("failed to read operation %d: %v", i+1, err)
		}
		if ops[i].i < 1 || ops[i].i > len(tc.array) || ops[i].j < 1 || ops[i].j > len(tc.array) {
			return fmt.Errorf("operation %d has invalid indices", i+1)
		}
	}
	arr := append([]int(nil), tc.array...)
	for _, op := range ops {
		arr[op.i-1] += arr[op.j-1]
		if arr[op.i-1] > 1_000_000_000 || arr[op.i-1] < -1_000_000_000 {
			return fmt.Errorf("values grew too large")
		}
	}
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return fmt.Errorf("array is not non-decreasing after operations")
		}
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([]int{2, 1}),
		makeCase([]int{1, 2, -10, 3}),
		makeCase([]int{-5, -3, -1}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(41) - 20
		}
		tests = append(tests, makeCase(arr))
	}
	return tests
}

func makeCase(arr []int) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(len(arr)))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{
		input: sb.String(),
		array: arr,
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
