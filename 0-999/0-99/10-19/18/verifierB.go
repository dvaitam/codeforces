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

func expected(n, d, m, l int64) int64 {
	limit := n * m
	x := d
	for {
		if x >= limit || x%m > l {
			return x
		}
		x += d
	}
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := int64(rng.Intn(20) + 1)
	d := int64(rng.Intn(20) + 1)
	m := int64(rng.Intn(20) + 1)
	l := int64(rng.Intn(int(m)))
	input := fmt.Sprintf("%d %d %d %d\n", n, d, m, l)
	return input, expected(n, d, m, l)
}

func runCase(bin, input string, exp int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
