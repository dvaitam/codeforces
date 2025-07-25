package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var primes []bool

func initPrimes() {
	n := 1000000
	primes = make([]bool, n+1)
	for i := 2; i <= n; i++ {
		primes[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if primes[p] {
			for m := p * p; m <= n; m += p {
				primes[m] = false
			}
		}
	}
}

func countPrimes(a, b int) int {
	c := 0
	for i := a; i <= b; i++ {
		if primes[i] {
			c++
		}
	}
	return c
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	initPrimes()
	for t := 0; t < 100; t++ {
		a := rand.Intn(1000000-1) + 2
		b := a + rand.Intn(1000000-a+1)
		in := fmt.Sprintf("%d %d\n", a, b)
		want := fmt.Sprintf("%d", countPrimes(a, b))
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != want {
			fmt.Printf("test %d failed: expected %q got %q\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
