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

const MOD = 1000000007

func sumDigitsBase(x, base int) int {
	sum := 0
	for x > 0 {
		sum += x % base
		x /= base
	}
	return sum
}

func expectedE(n int, k int) int {
	dp := 1
	last := make([]int, k)
	for j := 0; j < n; j++ {
		s := sumDigitsBase(j, k) % k
		ndp := (dp * 2) % MOD
		ndp = (ndp - last[s] + MOD) % MOD
		last[s] = dp
		dp = ndp
	}
	return dp
}

func genCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(30) + 1
	k := rng.Intn(5) + 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	exp := fmt.Sprintf("%d\n", expectedE(n, k))
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseE(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
