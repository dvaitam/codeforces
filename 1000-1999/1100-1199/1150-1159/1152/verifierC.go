package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) uint64 {
	g := gcd(a, b)
	return uint64(a/g) * uint64(b)
}

func expected(a, b int64) int64 {
	if a == b {
		return 0
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	bestK := int64(0)
	bestLCM := ^uint64(0)
	for d := int64(1); d*d <= diff; d++ {
		if diff%d != 0 {
			continue
		}
		for _, div := range []int64{d, diff / d} {
			if div == 0 {
				continue
			}
			k := (div - a%div) % div
			A := a + k
			B := b + k
			l := lcm(A, B)
			if l < bestLCM || (l == bestLCM && k < bestK) {
				bestLCM = l
				bestK = k
			}
		}
	}
	return bestK
}

func generateCase(rng *rand.Rand) (string, int64) {
	a := rng.Int63n(1_000_000) + 1
	b := rng.Int63n(1_000_000) + 1
	inp := fmt.Sprintf("%d %d\n", a, b)
	return inp, expected(a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	fixed := [][2]int64{
		{6, 10},
		{4, 4},
		{1, 9},
		{100000, 99999},
	}
	idx := 0
	for ; idx < len(fixed); idx++ {
		a := fixed[idx][0]
		b := fixed[idx][1]
		inp := fmt.Sprintf("%d %d\n", a, b)
		exp := strconv.FormatInt(expected(a, b), 10)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, exp, out, inp)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		inp, expVal := generateCase(rng)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strconv.FormatInt(expVal, 10) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", idx+1, expVal, out, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
