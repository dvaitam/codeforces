package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(x int64) int64 {
	bestK := int64(0)
	bestDiff := int64(1<<62 - 1)
	for k := int64(0); k < 4; k++ {
		angle := -x + 90*k
		mod := ((angle % 360) + 360) % 360
		if mod > 180 {
			mod -= 360
		}
		diff := int64(math.Abs(float64(mod)))
		if diff < bestDiff {
			bestDiff = diff
			bestK = k
		}
	}
	return bestK
}

func generateCase(rng *rand.Rand) (string, string) {
	x := int64(rng.Int63n(2_000_000_000_000_000_000) - 1_000_000_000_000_000_000)
	ans := solve(x)
	input := fmt.Sprintf("%d\n", x)
	expected := fmt.Sprintf("%d", ans)
	return input, expected
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
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierM.go /path/to/binary")
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
