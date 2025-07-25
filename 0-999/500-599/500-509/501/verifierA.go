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

func score(p, t int) int {
	term1 := 3 * p / 10
	term2 := p - (p/250)*t
	if term1 > term2 {
		return term1
	}
	return term2
}

func solveA(a, b, c, d int) string {
	s1 := score(a, c)
	s2 := score(b, d)
	switch {
	case s1 > s2:
		return "Misha"
	case s2 > s1:
		return "Vasya"
	default:
		return "Tie"
	}
}

func genCase(rng *rand.Rand) (string, string) {
	a := (rng.Intn(14) + 1) * 250
	b := (rng.Intn(14) + 1) * 250
	c := rng.Intn(181)
	d := rng.Intn(181)
	input := fmt.Sprintf("%d %d %d %d\n", a, b, c, d)
	expect := solveA(a, b, c, d)
	return input, expect
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
