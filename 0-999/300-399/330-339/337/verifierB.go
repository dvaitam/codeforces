package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solve(a, b, c, d int) (int, int) {
	t1 := a * d
	t2 := b * c
	var p, q int
	if t1 <= t2 {
		p = b*c - a*d
		q = b * c
	} else {
		p = a*d - b*c
		q = a * d
	}
	if p == 0 {
		return 0, 1
	}
	g := gcd(p, q)
	return p / g, q / g
}

func generateTest(rng *rand.Rand) (string, string) {
	a := rng.Intn(1000) + 1
	b := rng.Intn(1000) + 1
	c := rng.Intn(1000) + 1
	d := rng.Intn(1000) + 1
	p, q := solve(a, b, c, d)
	inp := fmt.Sprintf("%d %d %d %d\n", a, b, c, d)
	return inp, fmt.Sprintf("%d/%d", p, q)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inp)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n%s", t, err, out.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
