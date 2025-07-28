package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	input  string
	output string
}

func solveCaseA(n int, a []int, m int, b []int) []string {
	maxA := 0
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
	}
	maxB := 0
	for _, v := range b {
		if v > maxB {
			maxB = v
		}
	}
	first := "Alice"
	if maxA < maxB {
		first = "Bob"
	}
	second := "Alice"
	if maxB >= maxA {
		second = "Bob"
	}
	return []string{first, second}
}

func generateTests() []TestCase {
	rand.Seed(1)
	tests := make([]TestCase, 0, 20)
	for t := 0; t < 20; t++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		a := make([]int, n)
		b := make([]int, m)
		for i := range a {
			a[i] = rand.Intn(50) + 1
		}
		for i := range b {
			b[i] = rand.Intn(50) + 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')

		outLines := solveCaseA(n, append([]int(nil), a...), m, append([]int(nil), b...))
		out := strings.Join(outLines, "\n") + "\n"
		tests = append(tests, TestCase{sb.String(), out})
	}
	return tests
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			continue
		}
		g := strings.TrimSpace(got)
		e := strings.TrimSpace(tc.output)
		if g != e {
			fmt.Printf("Test %d failed. Expected %q, got %q\n", i+1, e, g)
		} else {
			passed++
		}
	}
	fmt.Printf("%d/%d tests passed\n", passed, len(tests))
}
