package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	switch strings.ToLower(filepath.Ext(bin)) {
	case ".go":
		cmd = exec.Command("go", "run", bin)
	case ".py":
		cmd = exec.Command("python3", bin)
	default:
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// validate checks the candidate's answer for one test case.
// a is the original array; output is the candidate's full output for that case.
func validate(a []int, output string) error {
	n := len(a)
	// Use flat token parsing to be robust to whitespace/newline variations.
	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("could not parse k: %v", err)
	}
	if k < n/2 || k > n {
		return fmt.Errorf("k=%d out of range [%d, %d]", k, n/2, n)
	}
	if len(tokens) != 1+k {
		return fmt.Errorf("k=%d but got %d value tokens", k, len(tokens)-1)
	}
	vals := make([]int, k)
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(tokens[1+i])
		if err != nil || (v != 0 && v != 1) {
			return fmt.Errorf("invalid value %q at position %d", tokens[1+i], i+1)
		}
		vals[i] = v
	}
	// Check vals is a subsequence of a.
	j := 0
	for _, v := range vals {
		for j < n && a[j] != v {
			j++
		}
		if j >= n {
			return fmt.Errorf("output is not a subsequence of the input array")
		}
		j++
	}
	// Check alternating sum == 0.
	alt := 0
	for i, v := range vals {
		if i%2 == 0 {
			alt += v
		} else {
			alt -= v
		}
	}
	if alt != 0 {
		return fmt.Errorf("alternating sum is %d, not 0", alt)
	}
	return nil
}

func genCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(20)*2 + 2 // even, between 2 and 40
	a := make([]int, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		v := rng.Intn(2)
		a[i] = v
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, a := genCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validate(a, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ngot:\n%s\ninput:\n%s", i+1, err, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
