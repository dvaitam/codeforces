package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input string
	a, b  int64
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
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	reader := bufio.NewReader(strings.NewReader(out))
	line1, err := reader.ReadString('\n')
	if err != nil && len(line1) == 0 {
		return fmt.Errorf("empty output")
	}
	line1 = strings.TrimSpace(line1)
	if line1 == "-1" {
		if solveReference(tc.a, tc.b) != -1 {
			return fmt.Errorf("solution exists but got -1")
		}
		return nil
	}
	k, err := strconv.Atoi(line1)
	if err != nil || k < 0 || k > 100 {
		return fmt.Errorf("invalid k value")
	}
	var ops []int64
	if k > 0 {
		line2, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read operations line")
		}
		line2 = strings.TrimSpace(line2)
		if line2 == "" {
			return fmt.Errorf("missing operations")
		}
		parts := strings.Fields(line2)
		if len(parts) != k {
			return fmt.Errorf("expected %d operations but got %d", k, len(parts))
		}
		for _, p := range parts {
			val, err := strconv.ParseInt(p, 10, 64)
			if err != nil || val < 0 || val > tc.a {
				return fmt.Errorf("invalid operation value: %s", p)
			}
			ops = append(ops, val)
		}
	}
	cur := tc.a
	for _, x := range ops {
		if x < 0 || x > cur {
			return fmt.Errorf("operation x=%d out of range for current a=%d", x, cur)
		}
		cur ^= x
	}
	if cur != tc.b {
		return fmt.Errorf("final value %d does not match target %d", cur, tc.b)
	}
	return nil
}

func genTests() []testCase {
	tests := []testCase{
		makeCase(9, 6),
		makeCase(13, 13),
		makeCase(292, 929),
		makeCase(405, 400),
		makeCase(998, 244),
		makeCase(353, 353),
	}
	for i := 0; i < 100; i++ {
		a := int64(i*10 + 1)
		b := int64((i*7 + 3) % 300)
		tests = append(tests, makeCase(a, b))
	}
	return tests
}

func makeCase(a, b int64) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	return testCase{input: sb.String(), a: a, b: b}
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

func solveReference(a, b int64) int {
	if a == b {
		return 0
	}
	if msb(b) > msb(a) {
		return -1
	}
	s := (int64(1) << (msb(a) + 1)) - 1
	if a != s && b != s {
		return 2
	}
	return 1
}

func msb(x int64) int {
	return bits.Len64(uint64(x)) - 1
}
