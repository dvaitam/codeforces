package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type query struct {
	c, m, x int64
}

type test struct {
	input    string
	expected string
}

func solveQuery(c, m, x int64) int64 {
	total := c + m + x
	teams := total / 3
	if c < teams {
		teams = c
	}
	if m < teams {
		teams = m
	}
	return teams
}

func solveCase(qs []query) string {
	var sb strings.Builder
	for i, qu := range qs {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", solveQuery(qu.c, qu.m, qu.x)))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func formatInput(qs []query) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(qs)))
	for _, qu := range qs {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.c, qu.m, qu.x))
	}
	return sb.String()
}

func fixedTests() []test {
	cases := [][]query{
		{{c: 1, m: 1, x: 1}},
		{{c: 10, m: 1, x: 1}},
		{{c: 0, m: 0, x: 0}},
		{{c: 5, m: 5, x: 5}, {c: 1, m: 2, x: 10}, {c: 7, m: 3, x: 0}},
	}
	var tests []test
	for _, qs := range cases {
		tests = append(tests, test{
			input:    formatInput(qs),
			expected: solveCase(qs),
		})
	}
	return tests
}

func randomQueries(rng *rand.Rand, count int, maxVal int64) []query {
	qs := make([]query, count)
	for i := 0; i < count; i++ {
		qs[i] = query{
			c: rng.Int63n(maxVal + 1),
			m: rng.Int63n(maxVal + 1),
			x: rng.Int63n(maxVal + 1),
		}
	}
	return qs
}

func randomTests(rng *rand.Rand, numTests int, maxQ int, maxVal int64) []test {
	tests := make([]test, 0, numTests)
	for len(tests) < numTests {
		qCount := rng.Intn(maxQ) + 1
		qs := randomQueries(rng, qCount, maxVal)
		tests = append(tests, test{
			input:    formatInput(qs),
			expected: solveCase(qs),
		})
	}
	return tests
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(1221))
	tests := fixedTests()
	tests = append(tests, randomTests(rng, 40, 5, 20)...)
	tests = append(tests, randomTests(rng, 30, 20, 1_000_000)...)
	tests = append(tests, randomTests(rng, 10, 10000, 100_000_000)...)
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
