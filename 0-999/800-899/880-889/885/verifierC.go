package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateTests() []int {
	r := rand.New(rand.NewSource(44))
	tests := make([]int, 100)
	for i := 0; i < 100; i++ {
		tests[i] = r.Intn(11) // 0..10
	}
	return tests
}

func expected(n int) int {
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		exp := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var ans int
		if _, err := fmt.Sscan(got, &ans); err != nil {
			fmt.Printf("Test %d: cannot parse output %q\n", i+1, got)
			os.Exit(1)
		}
		if ans != exp {
			fmt.Printf("Test %d failed. Input: %sExpected %d got %d\n", i+1, input, exp, ans)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
