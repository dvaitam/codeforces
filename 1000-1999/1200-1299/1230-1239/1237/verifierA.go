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

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) (string, []int) {
	n := r.Intn(20) + 2
	vals := make([]int, n)
	sum := 0
	for i := 0; i < n-1; i++ {
		vals[i] = r.Intn(2001) - 1000
		sum += vals[i]
	}
	vals[n-1] = -sum
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, v := range vals {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	return sb.String(), vals
}

// isValid checks that b is either floor(a/2) or ceil(a/2).
func isValid(a, b int) bool {
	if a%2 == 0 {
		return b == a/2
	}
	// a is odd: 2b must be a±1
	diff := 2*b - a
	return diff == 1 || diff == -1
}

func check(output string, as []int) error {
	r := strings.NewReader(output)
	n := len(as)
	bs := make([]int, n)
	for i := range bs {
		if _, err := fmt.Fscan(r, &bs[i]); err != nil {
			return fmt.Errorf("failed to parse output value %d: %v", i+1, err)
		}
	}
	sum := 0
	for i, b := range bs {
		if !isValid(as[i], b) {
			return fmt.Errorf("b[%d]=%d is not floor or ceil of %d/2", i+1, b, as[i])
		}
		sum += b
	}
	if sum != 0 {
		return fmt.Errorf("sum of b_i = %d, want 0", sum)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, as := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if err := check(got, as); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ngot:\n%s\ninput:\n%s", i, err, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
