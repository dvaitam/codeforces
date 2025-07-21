package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
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

func nextPrime(n int) int {
	for p := n + 1; ; p++ {
		if isPrime(p) {
			return p
		}
	}
}

func runCase(bin string, n, m int, expect string) error {
	input := fmt.Sprintf("%d %d\n", n, m)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %s got %s", expect, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}
	total := 0
	for _, n := range primes {
		np := nextPrime(n)
		for m := n + 1; m <= 50; m++ {
			exp := "NO"
			if m == np {
				exp = "YES"
			}
			if err := runCase(bin, n, m, exp); err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d m=%d)\n", total+1, err, n, m)
				os.Exit(1)
			}
			total++
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}
