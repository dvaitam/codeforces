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

var primes = []int{5, 7, 11, 13, 17, 19, 23, 29, 31}

func solve(n int, r int) float64 {
	nn := float64(n)
	rr := float64(r)
	angle := math.Pi / nn
	angle2 := angle / 2
	zj := angle
	h := math.Sin(zj) * rr
	d := math.Cos(zj) * rr
	p := math.Pi/2 - zj - angle2
	db := math.Tan(p) * h
	s := (h*d - h*db) * nn
	return s
}

func generateCase(rng *rand.Rand) (string, string) {
	n := primes[rng.Intn(len(primes))]
	r := rng.Intn(1000) + 1
	ans := solve(n, r)
	input := fmt.Sprintf("%d %d\n", n, r)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierP.go /path/to/binary")
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
