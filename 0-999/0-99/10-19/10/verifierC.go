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

func solveC(n int64) int64 {
	tot := int64(0)
	for i := int64(1); i <= n; i++ {
		tot -= n / i
	}
	var c [9]int64
	for i := int64(0); i < 9; i++ {
		c[i] = n / 9
		if n%9 >= i {
			c[i]++
		}
	}
	c[0]--
	for i := int64(0); i < 9; i++ {
		for j := int64(0); j < 9; j++ {
			for k := int64(0); k < 9; k++ {
				if (i*j)%9 == k {
					tot += c[i] * c[j] * c[k]
				}
			}
		}
	}
	return tot
}

func generateCaseC(rng *rand.Rand) (string, int64) {
	n := rng.Int63n(1_000_000) + 1
	return fmt.Sprintf("%d\n", n), solveC(n)
}

func runCaseC(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		if err := runCaseC(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
