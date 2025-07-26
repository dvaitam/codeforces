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

func solve(n, k int, a []int64) int64 {
	bestSpan := int64(1<<62 - 1)
	bestI := 0
	for i := 0; i+k < n; i++ {
		span := a[i+k] - a[i]
		if span < bestSpan {
			bestSpan = span
			bestI = i
		}
	}
	return (a[bestI] + a[bestI+k]) / 2
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var input strings.Builder
	var output strings.Builder
	fmt.Fprintf(&input, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 2
		k := rng.Intn(n-1) + 1
		arr := make([]int64, n)
		val := rng.Int63n(1_000_000)
		for j := 0; j < n; j++ {
			val += int64(rng.Intn(10) + 1)
			arr[j] = val
		}
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", arr[j])
		}
		input.WriteByte('\n')
		fmt.Fprintf(&output, "%d\n", solve(n, k, arr))
	}
	return input.String(), output.String()
}

func runCase(bin string, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
