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

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(bytes.TrimSpace(out)), err
}

var primes []uint64

func init() {
	// simple primes up to 1000
	isPrime := make([]bool, 1001)
	for i := 2; i <= 1000; i++ {
		isPrime[i] = true
	}
	for p := 2; p*p <= 1000; p++ {
		if isPrime[p] {
			for j := p * p; j <= 1000; j += p {
				isPrime[j] = false
			}
		}
	}
	for i := 2; i <= 1000; i++ {
		if isPrime[i] {
			primes = append(primes, uint64(i))
		}
	}
}

func genNumber(rng *rand.Rand) uint64 {
	p := primes[rng.Intn(len(primes))]
	choice := rng.Intn(4)
	switch choice {
	case 0:
		return p * p // p^2 -> 3 divisors
	case 1:
		return p * p * p // p^3 -> 4 divisors
	case 2:
		q := primes[rng.Intn(len(primes))]
		for q == p {
			q = primes[rng.Intn(len(primes))]
		}
		return p * q // product of two primes -> 4 divisors
	default:
		return p * p * p * p // p^4 -> 5 divisors
	}
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", genNumber(rng)))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	ref := "./refD.bin"
	if err := exec.Command("go", "build", "-o", ref, "1033D.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
