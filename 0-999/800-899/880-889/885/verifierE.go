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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func generateTests() [][2]int {
	r := rand.New(rand.NewSource(46))
	tests := make([][2]int, 100)
	for i := 0; i < 100; i++ {
		tests[i][0] = r.Intn(1000) + 1 // 1..1000
		tests[i][1] = r.Intn(1000) + 1
	}
	return tests
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t[0], t[1])
		exp := gcd(t[0], t[1])
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
