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

type testA struct {
	n, a, b, c int64
}

func solveA(n, a, b, c int64) int64 {
	var ans int64
	if n >= b && b-c <= a {
		k := (n-b)/(b-c) + 1
		ans += k
		n -= k * (b - c)
	}
	ans += n / a
	return ans
}

func genTests() []testA {
	rand.Seed(1)
	tests := make([]testA, 100)
	for i := range tests {
		n := rand.Int63n(1e12) + 1
		a := rand.Int63n(1e12) + 1
		b := rand.Int63n(1e12) + 1
		if b == 1 {
			b = 2
		}
		c := rand.Int63n(b-1) + 1
		tests[i] = testA{n: n, a: a, b: b, c: c}
	}
	// add some edge cases
	tests = append(tests, testA{n: 1, a: 1, b: 2, c: 1})
	tests = append(tests, testA{n: 100, a: 5, b: 10, c: 1})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n%d\n%d\n%d\n", t.n, t.a, t.b, t.c)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := fmt.Sprintf("%d", solveA(t.n, t.a, t.b, t.c))
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %s\nGot: %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
		if i%10 == 0 {
			time.Sleep(10 * time.Millisecond) // small delay to avoid runaway
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
