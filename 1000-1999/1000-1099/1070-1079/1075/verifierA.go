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
	n, x, y int64
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func expected(t Test) string {
	if t.x == t.n && t.y == t.n {
		return "Black"
	}
	left := max(abs(t.x-1), abs(t.y-1))
	right := max(abs(t.n-t.x), abs(t.n-t.y))
	if left > right {
		return "Black"
	}
	return "White"
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(0)
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		n := rand.Int63n(1_000_000) + 2
		x := rand.Int63n(n) + 1
		y := rand.Int63n(n) + 1
		tests = append(tests, Test{n, x, y})
	}
	tests = append(tests,
		Test{2, 1, 1},
		Test{2, 2, 2},
		Test{2, 2, 1},
		Test{10, 1, 10},
		Test{10, 10, 1},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.x, tc.y)
		exp := expected(tc)
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != exp {
			fmt.Printf("Test %d failed\nInput:%sExpected:%s\nGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
