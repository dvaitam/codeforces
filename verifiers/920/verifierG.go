package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func primeFactors(n int) []int {
	f := []int{}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			f = append(f, i)
			for n%i == 0 {
				n /= i
			}
		}
	}
	if n > 1 {
		f = append(f, n)
	}
	return f
}

func countCoprime(limit int64, factors []int) int64 {
	if limit <= 0 {
		return 0
	}
	m := len(factors)
	var bad int64
	for mask := 1; mask < (1 << m); mask++ {
		prod := int64(1)
		bits := 0
		for i := 0; i < m; i++ {
			if mask&(1<<i) != 0 {
				prod *= int64(factors[i])
				bits++
			}
		}
		if bits%2 == 1 {
			bad += limit / prod
		} else {
			bad -= limit / prod
		}
	}
	return limit - bad
}

func kth(x, p, k int) int64 {
	f := primeFactors(p)
	base := countCoprime(int64(x), f)
	low := int64(x) + 1
	high := int64(x) + int64(k*p) + 100
	for low < high {
		mid := (low + high) / 2
		if countCoprime(mid, f)-base >= int64(k) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(7)
	const tests = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, tests)
	expected := make([]int64, tests)
	for i := 0; i < tests; i++ {
		x := rand.Intn(1000) + 1
		p := rand.Intn(1000) + 1
		k := rand.Intn(20) + 1
		expected[i] = kth(x, p, k)
		fmt.Fprintf(&input, "%d %d %d\n", x, p, k)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to run binary:", err)
		os.Exit(1)
	}
	parts := strings.Fields(string(out))
	if len(parts) != tests {
		fmt.Printf("Expected %d numbers, got %d\n", tests, len(parts))
		os.Exit(1)
	}
	for i := 0; i < tests; i++ {
		val, err := strconv.ParseInt(parts[i], 10, 64)
		if err != nil || val != expected[i] {
			fmt.Printf("Test %d failed: expected %d got %s\n", i+1, expected[i], parts[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
