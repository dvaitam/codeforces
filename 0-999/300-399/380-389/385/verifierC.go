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

func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if isPrime[p] {
			for q := p * p; q <= n; q += p {
				isPrime[q] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func expected(xs []int, queries [][2]int) []int64 {
	maxR := 0
	for _, q := range queries {
		if q[1] > maxR {
			maxR = q[1]
		}
	}
	primes := sieve(maxR)
	results := make([]int64, len(queries))
	for idx, q := range queries {
		l, r := q[0], q[1]
		var sum int64
		for _, p := range primes {
			if p < l {
				continue
			}
			if p > r {
				break
			}
			cnt := 0
			for _, v := range xs {
				if v%p == 0 {
					cnt++
				}
			}
			sum += int64(cnt)
		}
		results[idx] = sum
	}
	return results
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(20) + 1
	xs := make([]int, n)
	for i := range xs {
		xs[i] = rng.Intn(100) + 2
	}
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(50) + 2
		r := l + rng.Intn(50)
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for _, qr := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qr[0], qr[1]))
	}
	return sb.String(), expected(xs, queries)
}

func runCase(bin, input string, expected []int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(expected) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		var val int64
		if _, err := fmt.Sscan(f, &val); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if val != expected[i] {
			return fmt.Errorf("expected %v got %v", expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
