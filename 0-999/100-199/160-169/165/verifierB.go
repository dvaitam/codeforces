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

func total(v, k int64) int64 {
	var sum int64
	for v > 0 {
		sum += v
		v /= k
	}
	return sum
}

func expected(n, k int64) int64 {
	low, high := int64(1), n
	for low < high {
		mid := (low + high) / 2
		if total(mid, k) >= n {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Int63n(1_000_000_000) + 1
	k := rng.Int63n(9) + 2 // 2..10
	if rng.Float64() < 0.2 {
		n = rng.Int63n(1000) + 1
	}
	input := fmt.Sprintf("%d %d\n", n, k)
	return input, expected(n, k)
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
