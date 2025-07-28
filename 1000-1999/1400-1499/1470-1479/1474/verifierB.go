package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Test struct {
	d        int
	expected int
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func nextPrime(x int) int {
	for {
		if isPrime(x) {
			return x
		}
		x++
	}
}

func genTests() []Test {
	r := rand.New(rand.NewSource(42))
	base := []int{1, 2, 3, 4, 5, 10, 20, 50, 100, 9999}
	tests := make([]Test, 0, 100)
	for _, d := range base {
		p1 := nextPrime(1 + d)
		p2 := nextPrime(p1 + d)
		tests = append(tests, Test{d: d, expected: p1 * p2})
	}
	for len(tests) < 100 {
		d := r.Intn(10000) + 1
		p1 := nextPrime(1 + d)
		p2 := nextPrime(p1 + d)
		tests = append(tests, Test{d: d, expected: p1 * p2})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t.d)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("time limit exceeded")
		return
	}
	if err != nil {
		fmt.Println("execution error:", err)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for i, t := range tests {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			return
		}
		var ans int
		_, err := fmt.Sscan(scanner.Text(), &ans)
		if err != nil || ans != t.expected {
			fmt.Printf("wrong answer on case %d: expected %d got %s\n", i+1, t.expected, scanner.Text())
			return
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		return
	}
	fmt.Println("OK")
}
