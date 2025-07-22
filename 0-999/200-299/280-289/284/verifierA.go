package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// totient computes Euler's totient of n
func totient(n int) int {
	result := n
	m := n
	for i := 2; i*i <= m; i++ {
		if m%i == 0 {
			for m%i == 0 {
				m /= i
			}
			result = result / i * (i - 1)
		}
	}
	if m > 1 {
		result = result / m * (m - 1)
	}
	return result
}

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// collect primes < 2000
	primes := []int{}
	for i := 2; i < 2000; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}

	for idx, p := range primes {
		input := fmt.Sprintf("%d\n", p)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (p=%d) runtime error: %v\n%s", idx+1, p, err, out.String())
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (p=%d) bad output: %v\n", idx+1, p, err)
			os.Exit(1)
		}
		expected := totient(p - 1)
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d (p=%d) expected %d got %d\n", idx+1, p, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
