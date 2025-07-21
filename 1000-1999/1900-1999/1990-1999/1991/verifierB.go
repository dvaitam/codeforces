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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2 // 2..10
	b := make([]int, n-1)
	for i := range b {
		b[i] = rng.Intn(16) // small numbers
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	// compute expected
	possible := true
	for i := 1; i < n-1 && possible; i++ {
		if ((b[i-1] & b[i+1]) &^ b[i]) != 0 {
			possible = false
		}
	}
	if !possible {
		return sb.String(), "-1\n"
	}
	a := make([]int, n)
	a[0] = b[0]
	for i := 1; i < n-1; i++ {
		a[i] = b[i-1] | b[i]
	}
	a[n-1] = b[n-2]
	var out strings.Builder
	for i, v := range a {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(strconv.Itoa(v))
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
