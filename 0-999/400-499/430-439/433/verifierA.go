package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input  string
	output string
}

func solve(weights []int) string {
	count100 := 0
	count200 := 0
	for _, w := range weights {
		if w == 100 {
			count100++
		} else {
			count200++
		}
	}
	total := count100*100 + count200*200
	if total%200 != 0 {
		return "NO"
	}
	if count100 == 0 && count200%2 == 1 {
		return "NO"
	}
	return "YES"
}

func generateTests() []test {
	rand.Seed(1)
	var tests []test
	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		weights := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				weights[j] = 100
			} else {
				weights[j] = 200
			}
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", weights[j])
		}
		sb.WriteByte('\n')
		out := solve(weights)
		tests = append(tests, test{sb.String(), out})
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.output {
			fmt.Printf("Test %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
