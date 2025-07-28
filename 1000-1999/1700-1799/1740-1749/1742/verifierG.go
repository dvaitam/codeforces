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
	arr []int
}

func generateCases() []testCase {
	rand.Seed(7)
	cases := make([]testCase, 100)
	for i := range cases {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(100)
		}
		cases[i] = testCase{arr: arr}
	}
	return cases
}

func expected(tc testCase) string {
	arr := append([]int(nil), tc.arr...)
	cur := 0
	limit := len(arr)
	if limit > 32 {
		limit = 32
	}
	for i := 0; i < limit; i++ {
		best := cur
		bestIdx := -1
		for j := i; j < len(arr); j++ {
			if cur|arr[j] > best {
				best = cur | arr[j]
				bestIdx = j
			}
		}
		if bestIdx != -1 {
			arr[i], arr[bestIdx] = arr[bestIdx], arr[i]
		}
		cur = best
	}
	var out strings.Builder
	for i, v := range arr {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(&out, v)
	}
	return out.String()
}

func buildIO(cases []testCase) (string, string) {
	var inBuilder strings.Builder
	var outBuilder strings.Builder
	fmt.Fprintf(&inBuilder, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&inBuilder, "%d\n", len(tc.arr))
		for i, v := range tc.arr {
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
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
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
