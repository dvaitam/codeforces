package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildCase(values []int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(values)))
	sb.WriteByte('\n')
	hasOne := false
	for i, v := range values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
		if v == 1 {
			hasOne = true
		}
	}
	sb.WriteByte('\n')
	expect := "1"
	if hasOne {
		expect = "-1"
	}
	return testCase{input: sb.String(), expected: expect}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	used := make(map[int]bool)
	values := make([]int, 0, n)
	for len(values) < n {
		v := rng.Intn(1000) + 1
		if !used[v] {
			used[v] = true
			values = append(values, v)
		}
	}
	return buildCase(values)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// deterministic cases
	cases = append(cases, buildCase([]int{1}))
	cases = append(cases, buildCase([]int{2}))
	cases = append(cases, buildCase([]int{1, 2, 3}))
	cases = append(cases, buildCase([]int{5, 6, 7}))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
