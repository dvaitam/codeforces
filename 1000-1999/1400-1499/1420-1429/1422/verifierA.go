package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	r := strings.NewReader(input)
	var t int
	fmt.Fscan(r, &t)
	var out strings.Builder
	for i := 0; i < t; i++ {
		var a, b, c int64
		fmt.Fscan(r, &a, &b, &c)
		fmt.Fprintf(&out, "%d\n", a+b+c-1)
	}
	return strings.TrimRight(out.String(), "\n")
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	tests := []test{}
	fixed := []string{
		"1\n1 2 3\n",
		"1\n10 20 30\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		var sb strings.Builder
		sb.WriteString("1\n")
		a := rng.Int63n(1e9) + 1
		b := rng.Int63n(1e9) + 1
		c := rng.Int63n(1e9) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", a, b, c)
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
