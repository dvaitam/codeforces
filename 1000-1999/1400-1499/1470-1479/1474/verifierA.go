package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	n        int
	b        string
	expected string
}

func computeA(b string) string {
	prev := -1
	a := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		d := int(b[i] - '0')
		if d+1 != prev {
			a[i] = '1'
			prev = d + 1
		} else {
			a[i] = '0'
			prev = d
		}
	}
	return string(a)
}

func genTests() []Test {
	r := rand.New(rand.NewSource(42))
	base := []string{"0", "1", "00", "01", "10", "11", "101", "010", "1111"}
	tests := make([]Test, 0, 100)
	for _, b := range base {
		tests = append(tests, Test{n: len(b), b: b, expected: computeA(b)})
	}
	for len(tests) < 100 {
		n := r.Intn(20) + 1
		bs := make([]byte, n)
		for i := 0; i < n; i++ {
			if r.Intn(2) == 1 {
				bs[i] = '1'
			} else {
				bs[i] = '0'
			}
		}
		b := string(bs)
		tests = append(tests, Test{n: n, b: b, expected: computeA(b)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t.n)
		fmt.Fprintln(&input, t.b)
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
		ans := strings.TrimSpace(scanner.Text())
		if ans != t.expected {
			fmt.Printf("wrong answer on case %d: expected %s got %s\n", i+1, t.expected, ans)
			return
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		return
	}
	fmt.Println("OK")
}
