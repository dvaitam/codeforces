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

func solve(k int) int { return k - 1 }

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []int{2, 3, 4, 5, 10, 100, 1000000000}
	for len(tests) < 100 {
		k := rng.Intn(1_000_000_000-1) + 2
		tests = append(tests, k)
	}
	for i, k := range tests {
		input := fmt.Sprintf("1\n%d\n", k)
		expected := solve(k)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != 1 {
			fmt.Fprintf(os.Stderr, "test %d: expected single integer got %q\n", i+1, out)
			os.Exit(1)
		}
		val, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse integer\n", i+1)
			os.Exit(1)
		}
		if val != expected {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d (k=%d)\n", i+1, expected, val, k)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
