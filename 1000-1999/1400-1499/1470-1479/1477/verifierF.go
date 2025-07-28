package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseF struct {
	n int
	k int
	l []int
}

func generateTests() []testCaseF {
	rand.Seed(42)
	tests := make([]testCaseF, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		k := rand.Intn(10) + 1
		l := make([]int, n)
		for j := 0; j < n; j++ {
			l[j] = rand.Intn(10) + 1
		}
		tests[i] = testCaseF{n: n, k: k, l: l}
	}
	return tests
}

func buildInput(t testCaseF) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", t.n, t.k)
	for i := 0; i < t.n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", t.l[i])
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := buildInput(t)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "execution failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(&out)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "no output on test %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		if line != "0" {
			fmt.Fprintf(os.Stderr, "test %d failed: expected 0 got %s\n", i+1, line)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
