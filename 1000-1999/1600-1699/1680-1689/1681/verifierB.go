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

func solveCaseB(n int, a []int, ops []int) int {
	shift := 0
	for _, b := range ops {
		shift = (shift + b) % n
	}
	return a[shift]
}

func generateTests() []TestCase {
	rand.Seed(2)
	tests := make([]TestCase, 0, 20)
	for t := 0; t < 20; t++ {
		n := rand.Intn(8) + 2
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(100) + 1
		}
		m := rand.Intn(5) + 1
		ops := make([]int, m)
		for i := range ops {
			ops[i] = rand.Intn(n)
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
		for i, v := range ops {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')

		out := fmt.Sprintf("%d\n", solveCaseB(n, append([]int(nil), a...), append([]int(nil), ops...)))
		tests = append(tests, TestCase{sb.String(), out})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, tc := range tests {
		got, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			continue
		}
		g := strings.TrimSpace(got)
		e := strings.TrimSpace(tc.output)
		if g != e {
			fmt.Printf("Test %d failed. Expected %s got %s\n", i+1, e, g)
		} else {
			passed++
		}
	}
	fmt.Printf("%d/%d tests passed\n", passed, len(tests))
}
