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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Generate a valid test case by constructing a permutation first, then deriving a[].
func genCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(10) + 2
	// Build a random permutation of 0..n-1
	perm := rng.Perm(n)
	// Derive a[] from the permutation: a[i] = perm[i] XOR perm[i+1]
	arr := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		arr[i] = perm[i] ^ perm[i+1]
	}
	return n, arr
}

// Validate that the output is a valid permutation of 0..n-1 with b[i]^b[i+1]=a[i].
func validate(n int, arr []int, output string) error {
	fields := strings.Fields(output)
	if len(fields) != n {
		return fmt.Errorf("expected %d values, got %d", n, len(fields))
	}
	b := make([]int, n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q at position %d", f, i)
		}
		b[i] = val
	}
	// Check permutation of 0..n-1
	seen := make([]bool, n)
	for i, v := range b {
		if v < 0 || v >= n {
			return fmt.Errorf("value %d at position %d out of range [0, %d)", v, i, n)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}
	// Check XOR property
	for i := 0; i < n-1; i++ {
		if b[i]^b[i+1] != arr[i] {
			return fmt.Errorf("b[%d]^b[%d] = %d, expected a[%d] = %d", i, i+1, b[i]^b[i+1], i, arr[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validate(n, arr, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput: %s\n", i+1, err, sb.String(), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
