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

type counts []int

func computeCounts(s string) counts {
	n := len(s)
	c := make(counts, n)
	for i := 0; i < n; i++ {
		run := 1
		for j := i + 1; j < n && s[j] == s[i]; j++ {
			run++
		}
		for l := 1; l <= run; l++ {
			c[l-1]++
		}
	}
	return c
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, counts) {
	n := rng.Intn(8) + 1
	if rng.Float64() < 0.1 {
		n = rng.Intn(12) + 1
	}
	letters := []byte{'a', 'b', 'c'}
	sb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sb)
	c := computeCounts(s)
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", c[i])
	}
	b.WriteByte('\n')
	return b.String(), c
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if len(got) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: length mismatch expected %d got %d\ninput:\n%s", i+1, len(exp), len(got), input)
			os.Exit(1)
		}
		cGot := computeCounts(got)
		for j := range exp {
			if j >= len(cGot) || cGot[j] != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: counts mismatch\ninput:\n%s", i+1, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
