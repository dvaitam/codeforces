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
	rand.Seed(4)
	cases := make([]testCase, 100)
	for i := range cases {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(1000) + 1
		}
		cases[i] = testCase{arr: arr}
	}
	return cases
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(tc testCase) string {
	pos := make([]int, 1001)
	for i, v := range tc.arr {
		if i+1 > pos[v] {
			pos[v] = i + 1
		}
	}
	ans := -1
	for i := 1; i <= 1000; i++ {
		if pos[i] == 0 {
			continue
		}
		for j := 1; j <= 1000; j++ {
			if pos[j] == 0 {
				continue
			}
			if gcd(i, j) == 1 {
				sum := pos[i] + pos[j]
				if sum > ans {
					ans = sum
				}
			}
		}
	}
	return fmt.Sprint(ans)
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
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
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
