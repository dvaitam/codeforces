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

type testCase struct {
	n uint64
}

func expected(n uint64) uint64 {
	m := n
	var p uint64
	for m%3 == 0 {
		m /= 3
		p++
	}
	denom := uint64(1)
	for i := uint64(0); i <= p; i++ {
		denom *= 3
	}
	return (n + denom - 1) / denom
}

func runCase(bin string, n uint64) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got uint64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := expected(n)
	if got != exp {
		return fmt.Errorf("for n=%d expected %d got %d", n, exp, got)
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

	cases := []uint64{1, 2, 3, 4, 27, 28, 81}
	for i := 0; i < 100; i++ {
		v := uint64(rng.Int63n(1e12) + 1)
		cases = append(cases, v)
	}

	for i, n := range cases {
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
