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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	exp := strings.ToLower(strings.TrimSpace(expect))
	act := strings.ToLower(strings.TrimSpace(actual))
	if exp != act {
		return fmt.Errorf("expected %s but got %s", expect, actual)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase(6, 4, [][]int{{2, 3}, {4, 2}, {6, 4}}),
		makeCase(1, 1, [][]int{{1, 1}}),
		makeCase(3, 2, [][]int{{1, 1}, {2, 2}, {3, 2}}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		k := rand.Intn(5) + 1
		tokens := make([][]int, k)
		for j := 0; j < k; j++ {
			tokens[j] = []int{rand.Intn(n) + 1, rand.Intn(m) + 1}
		}
		tests = append(tests, makeCase(n, m, tokens))
	}
	return tests
}

func makeCase(n, m int, tokens [][]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d %d %d\n", n, m, len(tokens)))
	for _, t := range tokens {
		sb.WriteString(fmt.Sprintf("%d %d\n", t[0], t[1]))
	}
	return testCase{
		input:  sb.String(),
		expect: solveReference(n, m, tokens),
	}
}

func solveReference(n, m int, tokens [][]int) string {
	nim := 0
	for _, t := range tokens {
		row := t[0]
		col := t[1]
		nim ^= (row + col) % 2
	}
	if nim != 0 {
		return "Mimo"
	}
	return "Yuyu"
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
