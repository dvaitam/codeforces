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

func expected(n int) int {
	return n/4 + (n%4)/2
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(1999) + 2 // 2..2000
	if n%2 == 1 {
		n++
		if n > 2000 {
			n -= 2
		}
	}
	return fmt.Sprintf("1\n%d\n", n), expected(n)
}

func runCase(exe, input string, expectedAns int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expectedAns {
		return fmt.Errorf("expected %d got %d", expectedAns, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edge := []int{2, 4, 6, 8, 10, 1998, 2000}
	for i, n := range edge {
		in := fmt.Sprintf("1\n%d\n", n)
		if err := runCase(exe, in, expected(n)); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "random case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
