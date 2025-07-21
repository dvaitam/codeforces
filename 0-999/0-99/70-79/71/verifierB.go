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

func compute(n, k, t int) []int {
	S := t * n * k / 100
	full := S / k
	rem := S % k
	res := make([]int, n)
	for i := 0; i < n; i++ {
		switch {
		case i < full:
			res[i] = k
		case i == full:
			res[i] = rem
		default:
			res[i] = 0
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := rng.Intn(10) + 1
	t := rng.Intn(101)
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d %d\n", n, k, t)
	res := compute(n, k, t)
	var out strings.Builder
	for i, v := range res {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprintf("%d", v))
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
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
