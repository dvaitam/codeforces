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

type interval struct {
	a, b int
}

func expected(x0 int, ivs []interval) int {
	left := 0
	right := 1000000000
	for _, iv := range ivs {
		a, b := iv.a, iv.b
		if a > b {
			a, b = b, a
		}
		if a > left {
			left = a
		}
		if b < right {
			right = b
		}
	}
	if left > right {
		return -1
	}
	if x0 < left {
		return left - x0
	}
	if x0 > right {
		return x0 - right
	}
	return 0
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	x0 := rng.Intn(1001)
	ivs := make([]interval, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x0))
	for i := 0; i < n; i++ {
		a := rng.Intn(1001)
		b := rng.Intn(1001)
		for b == a {
			b = rng.Intn(1001)
		}
		ivs[i] = interval{a, b}
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	return sb.String(), expected(x0, ivs)
}

func runCase(bin, input string, expect int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
