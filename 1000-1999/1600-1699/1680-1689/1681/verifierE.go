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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func solveQuery(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func generateTests() []TestCase {
	rand.Seed(5)
	tests := make([]TestCase, 0, 20)
	for t := 0; t < 20; t++ {
		n := rand.Intn(3) + 2
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n-1; i++ {
			sb.WriteString("1 1 1 1\n")
		}
		m := rand.Intn(3) + 1
		sb.WriteString(fmt.Sprintf("%d\n", m))
		res := make([]int, m)
		for i := 0; i < m; i++ {
			x1 := rand.Intn(10)
			y1 := rand.Intn(10)
			x2 := rand.Intn(10)
			y2 := rand.Intn(10)
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", x1, y1, x2, y2))
			res[i] = solveQuery(x1, y1, x2, y2)
		}
		outBuilder := strings.Builder{}
		for _, r := range res {
			outBuilder.WriteString(fmt.Sprintf("%d\n", r))
		}
		tests = append(tests, TestCase{sb.String(), outBuilder.String()})
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
