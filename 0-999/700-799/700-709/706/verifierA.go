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

func generateCase(rng *rand.Rand) (string, float64) {
	a := rng.Intn(201) - 100
	b := rng.Intn(201) - 100
	n := rng.Intn(100) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n%d\n", a, b, n)
	minT := math.Inf(1)
	for i := 0; i < n; i++ {
		x := rng.Intn(201) - 100
		y := rng.Intn(201) - 100
		v := rng.Intn(100) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", x, y, v)
		dist := math.Hypot(float64(x-a), float64(y-b))
		t := dist / float64(v)
		if t < minT {
			minT = t
		}
	}
	return sb.String(), minT
}

func runCase(bin, input string, expected float64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6*math.Max(1, math.Abs(expected)) {
		return fmt.Errorf("expected %.6f got %.6f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
