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

func pyramidVolume(side float64, k int) float64 {
	r := side / 2 / math.Sin(math.Pi/float64(k))
	h := math.Sqrt(side*side - r*r)
	area := r * r * math.Sin(2*math.Pi/float64(k)) / 2 * float64(k)
	return area * h / 3
}

func solve(l3, l4, l5 float64) float64 {
	ret := pyramidVolume(l3, 3)
	ret += pyramidVolume(l4, 4)
	ret += pyramidVolume(l5, 5)
	return ret
}

func generateCase(rng *rand.Rand) (string, string) {
	l3 := rng.Float64()*999 + 1
	l4 := rng.Float64()*999 + 1
	l5 := rng.Float64()*999 + 1
	ans := solve(l3, l4, l5)
	input := fmt.Sprintf("%.0f %.0f %.0f\n", l3, l4, l5)
	expected := fmt.Sprintf("%.10f", ans)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierQ.go /path/to/binary")
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
