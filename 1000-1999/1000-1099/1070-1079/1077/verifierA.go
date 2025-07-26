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

func expected(a, b, k int64) int64 {
	l := k / 2
	r := k - l
	return r*a - l*b
}

func runCase(exe, input string, exp []int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers, got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		if v != exp[i] {
			return fmt.Errorf("mismatch at position %d: expected %d got %d", i+1, exp[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		q := rng.Intn(20) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", q)
		exp := make([]int64, q)
		for i := 0; i < q; i++ {
			a := rng.Int63n(1e9) + 1
			b := rng.Int63n(1e9) + 1
			k := rng.Int63n(1e9) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", a, b, k)
			exp[i] = expected(a, b, k)
		}
		input := sb.String()
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
