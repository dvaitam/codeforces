package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseC struct {
	n int
	a []string
	b []string
	c []string
}

func solveC(n int, a, b, c []string) (int, int, int) {
	counts := make(map[string]int)
	for i := 0; i < n; i++ {
		counts[a[i]]++
		counts[b[i]]++
		counts[c[i]]++
	}
	scores := [3]int{}
	for i := 0; i < n; i++ {
		if counts[a[i]] == 1 {
			scores[0] += 3
		} else if counts[a[i]] == 2 {
			scores[0]++
		}
		if counts[b[i]] == 1 {
			scores[1] += 3
		} else if counts[b[i]] == 2 {
			scores[1]++
		}
		if counts[c[i]] == 1 {
			scores[2] += 3
		} else if counts[c[i]] == 2 {
			scores[2]++
		}
	}
	return scores[0], scores[1], scores[2]
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

var words = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

func generateTests() []testCaseC {
	rand.Seed(44)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		a := make([]string, n)
		b := make([]string, n)
		c := make([]string, n)
		for j := 0; j < n; j++ {
			a[j] = words[rand.Intn(len(words))]
			b[j] = words[rand.Intn(len(words))]
			c[j] = words[rand.Intn(len(words))]
		}
		tests[i] = testCaseC{n, a, b, c}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(tc.a[j])
		}
		sb.WriteByte('\n')
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(tc.b[j])
		}
		sb.WriteByte('\n')
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(tc.c[j])
		}
		sb.WriteByte('\n')
		input := sb.String()
		x, y, z := solveC(tc.n, tc.a, tc.b, tc.c)
		expected := fmt.Sprintf("%d %d %d", x, y, z)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("test %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
