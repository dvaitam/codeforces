package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func stepsBinary(s string) int {
	x := new(big.Int)
	x.SetString(s, 2)
	one := big.NewInt(1)
	steps := 0
	for x.Cmp(one) != 0 {
		if x.Bit(0) == 1 {
			x.Add(x, one)
		} else {
			x.Rsh(x, 1)
		}
		steps++
	}
	return steps
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(200) + 1
	b := make([]byte, n)
	b[0] = '1'
	for i := 1; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	s := string(b)
	exp := stepsBinary(s)
	input := fmt.Sprintf("%s\n", s)
	out := fmt.Sprintf("%d\n", exp)
	return input, out
}

func runCase(exe, input, expected string) error {
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
		return fmt.Errorf("expected %s got %s", exp, outStr)
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
