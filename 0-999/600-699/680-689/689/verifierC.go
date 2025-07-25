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

func countWays(n int64) int64 {
	var res int64
	for k := int64(2); k*k*k <= n; k++ {
		res += n / (k * k * k)
	}
	return res
}

func solve(m int64) int64 {
	low, high := int64(1), int64(1e18)
	for low < high {
		mid := (low + high) / 2
		if countWays(mid) >= m {
			high = mid
		} else {
			low = mid + 1
		}
	}
	if countWays(low) == m {
		return low
	}
	return -1
}

func runCase(bin string, m int64) error {
	input := fmt.Sprintf("%d\n", m)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := solve(m)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) int64 {
	return rng.Int63n(1_000_000_000_000) + 1 // up to 1e12
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m := randomCase(rng)
		if err := runCase(bin, m); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
