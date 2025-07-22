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

func solve(data string) string {
	var n, k, l, c, d, p, nl, np int
	fmt.Sscan(data, &n, &k, &l, &c, &d, &p, &nl, &np)
	totalDrink := k * l
	totalLime := c * d
	toastsByDrink := totalDrink / nl
	toastsByLime := totalLime
	toastsBySalt := p / np
	maxToasts := toastsByDrink
	if toastsByLime < maxToasts {
		maxToasts = toastsByLime
	}
	if toastsBySalt < maxToasts {
		maxToasts = toastsBySalt
	}
	res := maxToasts / n
	return fmt.Sprintf("%d\n", res)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := rng.Intn(1000) + 1
	l := rng.Intn(1000) + 1
	c := rng.Intn(1000) + 1
	d := rng.Intn(1000) + 1
	p := rng.Intn(1000) + 1
	nl := rng.Intn(1000) + 1
	np := rng.Intn(1000) + 1
	input := fmt.Sprintf("%d %d %d %d %d %d %d %d\n", n, k, l, c, d, p, nl, np)
	expected := solve(input)
	return input, expected
}

func runCase(exe, input, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
