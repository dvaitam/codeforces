package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(100) + 1
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", v)
	}
	in.WriteByte('\n')

	sort.Ints(a)
	var out strings.Builder
	for i, v := range a {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", v)
	}
	out.WriteByte('\n')
	return in.String(), out.String()
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, buf.String())
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
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
