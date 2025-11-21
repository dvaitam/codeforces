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
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
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
	if strings.TrimSpace(expect) != strings.TrimSpace(actual) {
		return fmt.Errorf("expected %s but got %s", expect, actual)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(2, []int{1}, []int{1}),
		makeCase(4, []int{1, 2, 3}, []int{2}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 2
		path := make([]int, n)
		for j := 0; j < n; j++ {
			path[j] = j + 1
		}
		start := rand.Intn(n) + 1
		tests = append(tests, makeCase(n, path, []int{start}))
	}
	return tests
}

func makeCase(n int, path []int, starts []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(starts)))
	for i := 1; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", path[i-1], path[i]))
	}
	for i, s := range starts {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", s))
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		expect: solveReference(n, starts[0]),
	}
}

func solveReference(n, start int) string {
	distStart := max(start-1, n-start)
	if distStart%2 == 1 {
		return "Ron"
	}
	return "Hermione"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
