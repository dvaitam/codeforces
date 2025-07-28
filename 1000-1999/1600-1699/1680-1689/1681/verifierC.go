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

func isSorted(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}

func solveCaseC(n int, a, b []int) string {
	if isSorted(a) && isSorted(b) {
		return "0\n"
	}
	ops := make([][2]int, 0)
	for i := 0; i < n-1; i++ {
		minA, minB := a[i], b[i]
		for j := i + 1; j < n; j++ {
			if a[j] < minA {
				minA = a[j]
			}
			if b[j] < minB {
				minB = b[j]
			}
		}
		index := -1
		for j := i + 1; j < n; j++ {
			if a[j] == minA && b[j] == minB {
				index = j
				break
			}
		}
		if index != -1 {
			a[i], a[index] = a[index], a[i]
			b[i], b[index] = b[index], b[i]
			ops = append(ops, [2]int{i + 1, index + 1})
		}
	}
	if isSorted(a) && isSorted(b) {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(ops)))
		for _, p := range ops {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		return sb.String()
	}
	return "-1\n"
}

func generateTests() []TestCase {
	rand.Seed(3)
	tests := make([]TestCase, 0, 20)
	for t := 0; t < 20; t++ {
		n := rand.Intn(5) + 2
		a := make([]int, n)
		b := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(20)
			b[i] = rand.Intn(20)
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
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out := solveCaseC(n, append([]int(nil), a...), append([]int(nil), b...))
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
			fmt.Printf("Test %d failed. Expected %q got %q\n", i+1, e, g)
		} else {
			passed++
		}
	}
	fmt.Printf("%d/%d tests passed\n", passed, len(tests))
}
