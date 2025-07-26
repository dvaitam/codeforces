package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func finite(p, q, b int64) bool {
	g := gcd(p, q)
	q /= g
	g = gcd(q, b)
	for g > 1 {
		for q%g == 0 {
			q /= g
		}
		g = gcd(q, b)
	}
	return q == 1
}

func generateA(rng *rand.Rand) (string, string) {
	p := rng.Int63n(1e9)
	q := rng.Int63n(1e9) + 1
	b := rng.Int63n(1e9-1) + 2
	exp := "Infinite"
	if finite(p, q, b) {
		exp = "Finite"
	}
	input := fmt.Sprintf("1\n%d %d %d\n", p, q, b)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateA(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
