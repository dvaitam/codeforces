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

func expected(l, r int64) int64 {
	pow10 := [...]int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000}
	ans := int64(0)
	for L := 1; L < len(pow10); L++ {
		lo := l
		if lo < pow10[L-1] {
			lo = pow10[L-1]
		}
		hi := r
		if hi > pow10[L]-1 {
			hi = pow10[L] - 1
		}
		if lo > hi {
			continue
		}
		M := pow10[L] - 1
		c1 := M / 2
		c2 := c1 + 1
		candidates := []int64{lo, hi, c1, c2}
		for _, n := range candidates {
			if n < lo || n > hi {
				continue
			}
			prod := n * (M - n)
			if prod > ans {
				ans = prod
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	l := rng.Int63n(1000000000) + 1
	r := l + rng.Int63n(1000000000-l+1)
	input := fmt.Sprintf("%d %d\n", l, r)
	return input, expected(l, r)
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
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
	cases := []struct{ l, r int64 }{
		{1, 1}, {1, 1000000000}, {123456789, 987654321},
	}
	for _, c := range cases {
		input := fmt.Sprintf("%d %d\n", c.l, c.r)
		if err := runCase(bin, input, expected(c.l, c.r)); err != nil {
			fmt.Fprintf(os.Stderr, "special case failed: %v\n", err)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
