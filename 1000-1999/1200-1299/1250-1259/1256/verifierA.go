package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(a, b, n, s int64) string {
	use := s / n
	if use > a {
		use = a
	}
	if s-use*n <= b {
		return "YES"
	}
	return "NO"
}

func generate() (string, string) {
	const T = 100
	rand.Seed(1)
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	for i := 0; i < T; i++ {
		a := rand.Int63n(1e9) + 1
		b := rand.Int63n(1e9) + 1
		n := rand.Int63n(1e9) + 1
		s := rand.Int63n(1e9) + 1
		fmt.Fprintf(&in, "%d %d %d %d\n", a, b, n, s)
		fmt.Fprintln(&out, solveCase(a, b, n, s))
	}
	return in.String(), out.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	out, err := runCandidate(bin, in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Fprintln(os.Stderr, "wrong answer")
		fmt.Fprintln(os.Stderr, "expected:\n"+exp)
		fmt.Fprintln(os.Stderr, "got:\n"+out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
