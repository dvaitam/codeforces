package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func validTriple(a, b, c, n int) bool {
	if a <= 0 || b <= 0 || c <= 0 {
		return false
	}
	if a%3 == 0 || b%3 == 0 || c%3 == 0 {
		return false
	}
	return a+b+c == n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []int{3, 4, 5, 6, 7, 8, 9, 10, 1000000000}
	for len(tests) < 100 {
		n := rng.Intn(1_000_000_000-3+1) + 3
		tests = append(tests, n)
	}

	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != 3 {
			fmt.Fprintf(os.Stderr, "test %d: expected 3 integers got %q\n", i+1, out)
			os.Exit(1)
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		c, err3 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse integers\n", i+1)
			os.Exit(1)
		}
		if !validTriple(a, b, c, n) {
			fmt.Fprintf(os.Stderr, "test %d failed: n=%d output=%d %d %d\n", i+1, n, a, b, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
