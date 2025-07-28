package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expected(n, m, rb, cb, rd, cd int) int {
	r, c := rb, cb
	dr, dc := 1, 1
	steps := 0
	for {
		if r == rd || c == cd {
			return steps
		}
		if r+dr > n || r+dr < 1 {
			dr = -dr
		}
		if c+dc > m || c+dc < 1 {
			dc = -dc
		}
		r += dr
		c += dc
		steps++
	}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		rb := rng.Intn(n) + 1
		cb := rng.Intn(m) + 1
		rd := rng.Intn(n) + 1
		cd := rng.Intn(m) + 1
		input := fmt.Sprintf("1\n%d %d %d %d %d %d\n", n, m, rb, cb, rd, cd)
		exp := fmt.Sprintf("%d", expected(n, m, rb, cb, rd, cd))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
