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

func generateTests() [][2]int {
	r := rand.New(rand.NewSource(42))
	tests := make([][2]int, 100)
	for i := 0; i < 100; i++ {
		tests[i][0] = r.Intn(2001) - 1000 // -1000..1000
		tests[i][1] = r.Intn(2001) - 1000
	}
	// Ensure at least one test with a==0 to trigger runtime error in buggy solutions
	tests[0][0] = 0
	tests[0][1] = r.Intn(2001) - 1000
	return tests
}

func expected(a, b int) int {
	return a + b
}

func run(bin string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	// Seed random with current time to avoid deterministic patterns when used externally
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t[0], t[1])
		exp := expected(t[0], t[1])
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
