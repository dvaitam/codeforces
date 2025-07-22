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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedAnswer(x, y, a, b int64) int64 {
	g := gcd(x, y)
	lcm := x / g * y
	if lcm > b {
		return 0
	}
	high := b / lcm
	low := (a - 1) / lcm
	return high - low
}

func generateCase(rng *rand.Rand) (int64, int64, int64, int64) {
	x := int64(rng.Intn(1000) + 1)
	y := int64(rng.Intn(1000) + 1)
	a := rng.Int63n(2_000_000_000) + 1
	b := rng.Int63n(2_000_000_000) + 1
	if a > b {
		a, b = b, a
	}
	return x, y, a, b
}

func runCase(bin string, x, y, a, b int64) error {
	input := fmt.Sprintf("%d %d %d %d\n", x, y, a, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := strconv.FormatInt(expectedAnswer(x, y, a, b), 10)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		x, y, a, b := generateCase(rng)
		if err := runCase(bin, x, y, a, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d %d %d\n", i+1, err, x, y, a, b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
