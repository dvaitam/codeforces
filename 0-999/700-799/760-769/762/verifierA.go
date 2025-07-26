package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func divisors(n int64) []int64 {
	small := make([]int64, 0)
	big := make([]int64, 0)
	for i := int64(1); i*i <= n; i++ {
		if n%i == 0 {
			small = append(small, i)
			if i != n/i {
				big = append(big, n/i)
			}
		}
	}
	for i := len(big)/2 - 1; i >= 0; i-- {
		opposite := len(big) - 1 - i
		big[i], big[opposite] = big[opposite], big[i]
	}
	return append(small, big...)
}

func expected(n, k int64) int64 {
	divs := divisors(n)
	if k <= 0 || k > int64(len(divs)) {
		return -1
	}
	return divs[k-1]
}

func run(bin string, input string) (string, error) {
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
		n := rng.Int63n(1_000_000_000_000) + 1
		divs := divisors(n)
		var k int64
		if rng.Intn(2) == 0 {
			k = int64(rng.Intn(len(divs))) + 1
		} else {
			k = int64(len(divs)) + int64(rng.Intn(5)+1)
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		exp := fmt.Sprintf("%d", expected(n, k))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			fmt.Printf("input:\n%s\noutput:\n%s\n", input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
