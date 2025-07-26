package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input  string
	output string
}

func solveB(n, k int) string {
	var sb strings.Builder
	sb.WriteString("YES\n")
	blank := strings.Repeat(".", n)
	if k%2 == 0 {
		half := k / 2
		row := "." + strings.Repeat("#", half) + strings.Repeat(".", n-1-half)
		sb.WriteString(blank + "\n")
		sb.WriteString(row + "\n")
		sb.WriteString(row + "\n")
		sb.WriteString(blank + "\n")
	} else if k >= n-2 {
		sb.WriteString(blank + "\n")
		row2 := "." + strings.Repeat("#", n-2) + "."
		rem := k - (n - 2)
		half := rem / 2
		midDots := (n - 2) - rem
		row3 := "." + strings.Repeat("#", half) + strings.Repeat(".", midDots) + strings.Repeat("#", half) + "."
		sb.WriteString(row2 + "\n")
		sb.WriteString(row3 + "\n")
		sb.WriteString(blank + "\n")
	} else {
		sb.WriteString(blank + "\n")
		left := (n - k) / 2
		row2 := strings.Repeat(".", left) + strings.Repeat("#", k) + strings.Repeat(".", left)
		sb.WriteString(row2 + "\n")
		sb.WriteString(blank + "\n")
		sb.WriteString(blank + "\n")
	}
	return sb.String()
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(10)*2 + 3 // odd between 3 and 21
		k := rand.Intn(2*(n-2) + 1)
		input := fmt.Sprintf("%d %d\n", n, k)
		output := solveB(n, k)
		tests = append(tests, Test{input: input, output: output})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(binary, t.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Test %d failed. Input: %q\nExpected: %qGot: %q\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
