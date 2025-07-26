package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "1176D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func sieve(limit int) []int {
	isPrime := make([]bool, limit+1)
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if !isPrime[i] {
			primes = append(primes, i)
			for j := i * 2; j <= limit; j += i {
				isPrime[j] = true
			}
		}
	}
	return primes
}

func spf(limit int) []int {
	spf := make([]int, limit+1)
	for i := 2; i <= limit; i++ {
		if spf[i] == 0 {
			for j := i; j <= limit; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	return spf
}

func genCase(r *rand.Rand, primes []int, spf []int) string {
	n := r.Intn(50) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = r.Intn(200000-2) + 2
	}
	b := make([]int, 0, 2*n)
	for _, v := range a {
		b = append(b, v)
		if spf[v] == v {
			b = append(b, primes[v-1])
		} else {
			b = append(b, v/spf[v])
		}
	}
	// shuffle b
	r.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, x := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x))
	}
	sb.WriteByte('\n')
	return sb.String()
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	primes := sieve(2750131)
	spfArr := spf(200000)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng, primes, spfArr)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
