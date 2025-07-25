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

func solve(x1, y1, x2, y2 int64) int64 {
	width := x2 - x1 + 1
	height := y2 - y1 + 1
	total := width * height
	ans := total / 2
	if total%2 == 1 {
		if ((x1 ^ y1) & 1) == 0 {
			ans++
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	w := int64(rng.Intn(1000))*2 + 2
	h := int64(rng.Intn(1000)) + 1
	x1 := int64(rng.Intn(2000)) - 1000
	y1 := int64(rng.Intn(2000)) - 1000
	x2 := x1 + w - 1
	y2 := y1 + h - 1
	ans := solve(x1, y1, x2, y2)
	input := fmt.Sprintf("%d %d %d %d\n", x1, y1, x2, y2)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
