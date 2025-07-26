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

func linearSieve(n int) []int {
	spf := make([]int, n+1)
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > n {
				break
			}
			spf[i*p] = p
		}
	}
	return spf
}

func primeFactors(n int, spf []int) []int {
	res := make([]int, 0)
	for n > 1 {
		p := spf[n]
		res = append(res, p)
		for n%p == 0 {
			n /= p
		}
	}
	return res
}

var sieve = linearSieve(1000000)

func solve(x2 int) int {
	factors2 := primeFactors(x2, sieve)
	best := x2
	for _, p2 := range factors2 {
		start := x2 - p2 + 1
		if start < 3 {
			start = 3
		}
		for x1 := start; x1 <= x2; x1++ {
			if p2 >= x1 {
				continue
			}
			if (x2-x1)%p2 != 0 {
				continue
			}
			factors1 := primeFactors(x1, sieve)
			for _, p1 := range factors1 {
				if p1 >= x1 {
					continue
				}
				cand := x1 - p1 + 1
				if cand < p1+1 {
					cand = p1 + 1
				}
				if cand < 3 {
					cand = 3
				}
				if cand <= x1 && cand < best {
					best = cand
				}
			}
		}
	}
	return best
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		x2 := rng.Intn(1000-4) + 4
		// ensure composite
		for sieve[x2] == x2 {
			x2 = rng.Intn(1000-4) + 4
		}
		input := fmt.Sprintf("%d\n", x2)
		expected := fmt.Sprintf("%d", solve(x2))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
