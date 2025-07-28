package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	a       []int
	queries []int
}

func generateCases() []testCase {
	rand.Seed(5)
	cases := make([]testCase, 100)
	for i := range cases {
		n := rand.Intn(10) + 1
		q := rand.Intn(10) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(10) + 1
		}
		queries := make([]int, q)
		for j := range queries {
			queries[j] = rand.Intn(12)
		}
		cases[i] = testCase{a: a, queries: queries}
	}
	return cases
}

func expected(tc testCase) string {
	n := len(tc.a)
	prefMax := make([]int, n)
	prefSum := make([]int64, n)
	curMax := 0
	var curSum int64
	for i, v := range tc.a {
		if v > curMax {
			curMax = v
		}
		curSum += int64(v)
		prefMax[i] = curMax
		prefSum[i] = curSum
	}
	var out strings.Builder
	for idx, k := range tc.queries {
		pos := sort.Search(n, func(j int) bool { return prefMax[j] > k })
		var ans int64
		if pos > 0 {
			ans = prefSum[pos-1]
		}
		if idx > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(&out, ans)
	}
	return out.String()
}

func buildIO(cases []testCase) (string, string) {
	var inBuilder strings.Builder
	var outBuilder strings.Builder
	fmt.Fprintf(&inBuilder, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&inBuilder, "%d %d\n", len(tc.a), len(tc.queries))
		for i, v := range tc.a {
			if i > 0 {
				inBuilder.WriteByte(' ')
			}
			fmt.Fprint(&inBuilder, v)
		}
		inBuilder.WriteByte('\n')
		for i, v := range tc.queries {
			if i > 0 {
				inBuilder.WriteByte(' ')
			}
			fmt.Fprint(&inBuilder, v)
		}
		inBuilder.WriteByte('\n')
		fmt.Fprintln(&outBuilder, expected(tc))
	}
	return inBuilder.String(), outBuilder.String()
}

func run(binary, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(binary, ".go") {
		cmd = exec.Command("go", "run", binary)
	} else {
		cmd = exec.Command(binary)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func normalizeTokens(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return strings.Fields(s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	cases := generateCases()
	input, expectedOutput := buildIO(cases)
	actualOutput, err := run(binary, input)
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		os.Exit(1)
	}
	if strings.Join(normalizeTokens(actualOutput), " ") != strings.Join(normalizeTokens(expectedOutput), " ") {
		fmt.Println("Wrong answer")
		fmt.Println("Expected:")
		fmt.Println(expectedOutput)
		fmt.Println("Got:")
		fmt.Println(actualOutput)
		os.Exit(1)
	}
	fmt.Println("All test cases passed!")
}
