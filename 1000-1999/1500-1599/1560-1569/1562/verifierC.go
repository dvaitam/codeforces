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

func generateCase(rng *rand.Rand) (int, string) {
	n := rng.Intn(50) + 2
	sb := make([]byte, n)
	for i := range sb {
		if rng.Intn(2) == 0 {
			sb[i] = '0'
		} else {
			sb[i] = '1'
		}
	}
	// Occasionally generate all 1s
	if rng.Intn(10) == 0 {
		for i := range sb {
			sb[i] = '1'
		}
	}
	return n, string(sb)
}

func verify(n int, s string, outStr string) error {
	var l1, r1, l2, r2 int
	_, err := fmt.Sscan(outStr, &l1, &r1, &l2, &r2)
	if err != nil {
		return fmt.Errorf("failed to scan output: %v", err)
	}

	if l1 < 1 || l1 > n || r1 < 1 || r1 > n {
		return fmt.Errorf("range 1 out of bounds: %d %d (n=%d)", l1, r1, n)
	}
	if l2 < 1 || l2 > n || r2 < 1 || r2 > n {
		return fmt.Errorf("range 2 out of bounds: %d %d (n=%d)", l2, r2, n)
	}

	if l1 > r1 {
		return fmt.Errorf("l1 > r1")
	}
	if l2 > r2 {
		return fmt.Errorf("l2 > r2")
	}

	minLen := n / 2
	if r1-l1+1 < minLen {
		return fmt.Errorf("range 1 too short: len %d < %d", r1-l1+1, minLen)
	}
	if r2-l2+1 < minLen {
		return fmt.Errorf("range 2 too short: len %d < %d", r2-l2+1, minLen)
	}

	if l1 == l2 && r1 == r2 {
		return fmt.Errorf("ranges are identical")
	}

	sub1 := s[l1-1 : r1]
	sub2 := s[l2-1 : r2]

	val1 := new(big.Int)
	val1.SetString(sub1, 2)

	val2 := new(big.Int)
	val2.SetString(sub2, 2)

	if val2.Sign() == 0 {
		if val1.Sign() != 0 {
			return fmt.Errorf("f(t)=%s (non-zero) is not multiple of f(w)=0", val1.String())
		}
	} else {
		rem := new(big.Int)
		rem.Mod(val1, val2)
		if rem.Sign() != 0 {
			return fmt.Errorf("f(t)=%s is not multiple of f(w)=%s", val1.String(), val2.String())
		}
	}
	return nil
}

func runCase(exe string, n int, s string) error {
	input := fmt.Sprintf("1\n%d\n%s\n", n, s)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	outStr := strings.TrimSpace(out.String())
	if err := verify(n, s, outStr); err != nil {
		return fmt.Errorf("%v\nOutput: %s", err, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, s := generateCase(rng)
		if err := runCase(exe, n, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d\n%s\n", i+1, err, n, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}