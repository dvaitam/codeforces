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

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest() ([]byte, int, int) {
	n := rand.Intn(49) + 2 // 2..50
	k := rand.Intn(n)
	return []byte(fmt.Sprintf("1\n%d %d\n", n, k)), n, k
}

// validate checks that the output is a valid permutation of 1..n with exactly k inversions.
func validate(output string, n, k int) error {
	tokens := strings.Fields(output)
	if len(tokens) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens))
	}
	perm := make([]int, n)
	used := make([]bool, n+1)
	for i, tok := range tokens {
		var v int
		if _, err := fmt.Sscan(tok, &v); err != nil {
			return fmt.Errorf("bad token %q: %v", tok, err)
		}
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range [1,%d]", v, n)
		}
		if used[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		used[v] = true
		perm[i] = v
	}
	// Count inversions
	inv := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if perm[i] > perm[j] {
				inv++
			}
		}
	}
	if inv != k {
		return fmt.Errorf("expected %d inversions, got %d", k, inv)
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		return
	}
	cand := os.Args[1]

	for i := 1; i <= 100; i++ {
		in, n, k := genTest()
		got, err := run(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := validate(got, n, k); err != nil {
			fmt.Printf("wrong answer on test %d (n=%d, k=%d): %v\noutput: %s\n", i, n, k, err, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
