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

const sievelimit = 1000001

var primes []int64
var sieve [sievelimit]bool

func initSieve() {
	for i := 2; i < sievelimit; i++ {
		sieve[i] = true
	}
	for i := 2; i*i < sievelimit; i++ {
		if sieve[i] {
			for j := i * i; j < sievelimit; j += i {
				sieve[j] = false
			}
		}
	}
	for i := 2; i < sievelimit; i++ {
		if sieve[i] {
			primes = append(primes, int64(i))
		}
	}
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n < sievelimit {
		return sieve[n]
	}
	for _, p := range primes {
		if p*p > n {
			break
		}
		if n%p == 0 {
			return false
		}
	}
	return true
}

// solveD counts positive integers x with sigma(x) == A.
// sigma is multiplicative: sigma(p1^a1 * p2^a2 * ...) = product of sigma(pi^ai).
// sigma(p^k) = 1 + p + p^2 + ... + p^k.
// f(idx, a) counts n with sigma(n)==a and all prime factors of n >= primes[idx].
func solveD(A int64) int64 {
	var f func(idx int, a int64) int64
	f = func(idx int, a int64) int64 {
		if a == 1 {
			return 1
		}
		count := int64(0)
		// Try each sieve prime p >= primes[idx].
		// J(x) = sum of unitary divisors: sigma*(p^k) = p^k + 1.
		for i := idx; i < len(primes); i++ {
			p := primes[i]
			if p+1 > a {
				break // sigma*(p) = p+1 > a, no solution possible
			}
			// Try p^1, p^2, p^3, ...: sigma*(p^k) = p^k + 1
			pk := p      // p^k
			s := pk + 1  // sigma*(p^1) = p+1
			for s <= a {
				if a%s == 0 {
					count += f(i+1, a/s)
				}
				pk *= p
				if pk > a {
					break
				}
				s = pk + 1 // sigma*(p^{k+1}) = p^{k+1} + 1
			}
		}
		// Handle large prime: if a-1 is prime and > sieve limit,
		// sigma(a-1) = a, giving a valid x = a-1.
		// (Small primes a-1 are already covered by the sieve loop above.)
		candidate := a - 1
		if candidate >= sievelimit && isPrime(candidate) {
			count++
		}
		return count
	}
	return f(0, A)
}

func runExe(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testInput struct{ text string }

func buildTests() []testInput {
	tests := []testInput{
		{text: "1\n"},
		{text: "2\n"},
		{text: "3\n"},
		{text: "24\n"},
		{text: "1000000000000\n"},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		val := rng.Int63n(1_000_000_000_000) + 1
		tests = append(tests, testInput{text: fmt.Sprintf("%d\n", val)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]
	if strings.HasSuffix(candidate, ".go") {
		tmp, err := os.CreateTemp("", "verifierD-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("go", "build", "-o", tmp.Name(), candidate).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		candidate = tmp.Name()
	} else if strings.HasSuffix(candidate, ".c") {
		src, err := os.ReadFile(candidate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read source: %v\n", err)
			os.Exit(1)
		}
		// Fix Windows-specific format specifier and inline linkage for Linux.
		src = bytes.ReplaceAll(src, []byte("%I64d"), []byte("%lld"))
		src = bytes.ReplaceAll(src, []byte("inline "), []byte("static inline "))
		tmpSrc, err := os.CreateTemp("", "verifierD-csrc-*.c")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp source: %v\n", err)
			os.Exit(1)
		}
		tmpSrc.Write(src)
		tmpSrc.Close()
		defer os.Remove(tmpSrc.Name())
		tmpBin, err := os.CreateTemp("", "verifierD-cbin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp binary: %v\n", err)
			os.Exit(1)
		}
		tmpBin.Close()
		defer os.Remove(tmpBin.Name())
		out, err := exec.Command("gcc", "-O2", "-o", tmpBin.Name(), tmpSrc.Name()).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		candidate = tmpBin.Name()
	}

	initSieve()
	tests := buildTests()
	for idx, test := range tests {
		expected := fmt.Sprintf("%d", solveD(
			func() int64 {
				var v int64
				fmt.Sscan(test.text, &v)
				return v
			}(),
		))
		got, err := runExe(candidate, test.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput: %s\n", idx+1, err, test.text)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\ninput: %sexpected: %s\ngot: %s\n",
				idx+1, test.text, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
