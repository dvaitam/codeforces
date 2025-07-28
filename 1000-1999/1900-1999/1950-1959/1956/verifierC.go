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

func solveCase(n int) string {
	var out bytes.Buffer
	var sum int64
	for i := 1; i <= n; i++ {
		sum += int64(2*i-1) * int64(i)
	}
	fmt.Fprintf(&out, "%d %d\n", sum, 2*n-1)
	// first operation
	fmt.Fprintf(&out, "1 %d", n)
	for i := 1; i <= n; i++ {
		out.WriteByte(' ')
		fmt.Fprintf(&out, "%d", i)
	}
	out.WriteByte('\n')
	for i := 1; i < n; i++ {
		fmt.Fprintf(&out, "2 %d", n-i)
		for j := 1; j <= n; j++ {
			out.WriteByte(' ')
			fmt.Fprintf(&out, "%d", j)
		}
		out.WriteByte('\n')
		fmt.Fprintf(&out, "1 %d", n-i)
		for j := 1; j <= n; j++ {
			out.WriteByte(' ')
			fmt.Fprintf(&out, "%d", j)
		}
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(2) + 1
	var in bytes.Buffer
	var out bytes.Buffer
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(3) + 1
		fmt.Fprintf(&in, "%d\n", n)
		out.WriteString(solveCase(n))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return in.String(), strings.TrimSpace(out.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
