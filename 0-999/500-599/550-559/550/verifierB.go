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

func expected(n, l, r, x int, c []int) string {
	count := 0
	for mask := 0; mask < (1 << n); mask++ {
		if bitsOnes(mask) < 2 {
			continue
		}
		sum := 0
		minv := 1<<31 - 1
		maxv := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				v := c[i]
				sum += v
				if v < minv {
					minv = v
				}
				if v > maxv {
					maxv = v
				}
			}
		}
		if sum >= l && sum <= r && maxv-minv >= x {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

func bitsOnes(x int) int {
	c := 0
	for x > 0 {
		c += x & 1
		x >>= 1
	}
	return c
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	l := rng.Intn(50)
	r := l + rng.Intn(50) + 1
	x := rng.Intn(20) + 1
	c := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, l, r, x)
	for i := 0; i < n; i++ {
		c[i] = rng.Intn(40) + 1
		fmt.Fprintf(&sb, "%d ", c[i])
	}
	sb.WriteString("\n")
	input := sb.String()
	exp := expected(n, l, r, x, c)
	return input, exp
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
		return fmt.Errorf("expected %q got %q", expected, got)
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
