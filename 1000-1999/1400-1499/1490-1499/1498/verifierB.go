package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveCase(n int, W int64, weights []int) int {
	counts := make([]int, 31)
	for _, w := range weights {
		idx := bits.TrailingZeros(uint(w))
		counts[idx]++
	}
	remaining := n
	height := 0
	for remaining > 0 {
		height++
		remW := W
		for j := 30; j >= 0; j-- {
			if counts[j] == 0 {
				continue
			}
			width := int64(1 << j)
			if width > remW {
				continue
			}
			maxFit := int(remW / width)
			if maxFit > counts[j] {
				maxFit = counts[j]
			}
			counts[j] -= maxFit
			remW -= width * int64(maxFit)
			remaining -= maxFit
		}
	}
	return height
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(4) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		W := int64(rng.Intn(1_000_000_000) + 1)
		in.WriteString(fmt.Sprintf("%d %d\n", n, W))
		weights := make([]int, n)
		maxIdx := bits.Len64(uint64(W)) - 1
		for j := 0; j < n; j++ {
			idx := rng.Intn(maxIdx + 1)
			weights[j] = 1 << idx
			in.WriteString(fmt.Sprintf("%d", weights[j]))
			if j+1 < n {
				in.WriteByte(' ')
			} else {
				in.WriteByte('\n')
			}
		}
		out.WriteString(fmt.Sprintf("%d\n", solveCase(n, W, weights)))
	}
	return in.String(), out.String()
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
