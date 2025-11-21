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
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
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
	expectFields := strings.Fields(expect)
	actualFields := strings.Fields(actual)
	if len(expectFields) != len(actualFields) {
		return fmt.Errorf("expected %d bits but got %d", len(expectFields), len(actualFields))
	}
	for i := range expectFields {
		if expectFields[i] != actualFields[i] {
			return fmt.Errorf("mismatch at bit %d: expected %s got %s", i, expectFields[i], actualFields[i])
		}
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest([]int{0, 0, 0}),
		makeTest([]int{1, 1, 1}),
		makeTest([]int{1, 0, 1, 0}),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(10) + 1
		bits := make([]int, n)
		for j := 0; j < n; j++ {
			bits[j] = rand.Intn(2)
		}
		tests = append(tests, makeTest(bits))
	}
	return tests
}

func makeTest(bits []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(bits))
	for i, b := range bits {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b))
	}
	sb.WriteByte('\n')
	expect := increment(bits)
	return testCase{
		input:  sb.String(),
		expect: expect,
	}
}

func increment(bits []int) string {
	tmp := append([]int(nil), bits...)
	carry := 1
	for i := 0; i < len(tmp); i++ {
		if carry == 0 {
			break
		}
		if tmp[i] == 0 {
			tmp[i] = 1
			carry = 0
		} else {
			tmp[i] = 0
		}
	}
	var sb strings.Builder
	for i, b := range tmp {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b))
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
