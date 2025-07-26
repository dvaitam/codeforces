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
	n, A, B, C, T int
	times         []int
}

// runCandidate executes the given binary or Go source with the provided input.
func runCandidate(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func solve(tc testCase) int64 {
	if tc.C <= tc.B {
		return int64(tc.n) * int64(tc.A)
	}
	var sum int64
	for _, t := range tc.times {
		sum += int64(tc.T - t)
	}
	return int64(tc.n)*int64(tc.A) + int64(tc.C-tc.B)*sum
}

func generateRandomCase() testCase {
	n := rand.Intn(20) + 1 // keep runtime small
	A := rand.Intn(1000) + 1
	B := rand.Intn(1000) + 1
	C := rand.Intn(1000) + 1
	T := rand.Intn(1000) + 1
	times := make([]int, n)
	for i := range times {
		times[i] = rand.Intn(T) + 1
	}
	return testCase{n, A, B, C, T, times}
}

func generateTests() []testCase {
	rand.Seed(2)
	cases := make([]testCase, 0, 100)
	for i := 0; i < 97; i++ {
		cases = append(cases, generateRandomCase())
	}
	// edge case: all minimum values
	cases = append(cases, testCase{1, 1, 1, 1, 1, []int{1}})
	// edge case: large values
	large := testCase{n: 1000, A: 1000, B: 999, C: 1000, T: 1000, times: make([]int, 1000)}
	for i := range large.times {
		large.times[i] = 1
	}
	cases = append(cases, large)
	// another large case with times at T
	large2 := testCase{n: 1000, A: 1, B: 0, C: 1000, T: 1000, times: make([]int, 1000)}
	for i := range large2.times {
		large2.times[i] = 1000
	}
	cases = append(cases, large2)
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", tc.n, tc.A, tc.B, tc.C, tc.T))
		for j, v := range tc.times {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := strconv.FormatInt(solve(tc), 10)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, got)
			fmt.Printf("input:\n%s", input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expected, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
