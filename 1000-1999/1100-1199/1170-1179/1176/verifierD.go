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

// validateCandidate checks whether the candidate output represents a valid selection
// of n numbers from b according to the problem pairing rules.
func validateCandidate(out string, n int, b []int, primes []int, spf []int) error {
	toks := strings.Fields(strings.TrimSpace(out))
	if len(toks) != n {
		return fmt.Errorf("wrong number of outputs: expected %d got %d", n, len(toks))
	}
	// frequency map of b values
	freq := make(map[int]int, len(b))
	for _, v := range b {
		freq[v]++
	}
	// Helper to check small-prime using spf
	isSmallPrime := func(x int) bool {
		if x <= 1 || x >= len(spf) {
			return false
		}
		return spf[x] == x
	}
	// Process each output value
	for _, s := range toks {
		var v int
		fmt.Sscanf(s, "%d", &v)
		used := false
		// Try composite case first if available and consistent
		if v < len(spf) && v >= 2 && !isSmallPrime(v) {
			if freq[v] > 0 {
				y := v / spf[v]
				if freq[y] > 0 {
					freq[v]--
					freq[y]--
					used = true
				}
			}
		}
		if !used {
			// Try prime-index case
			if v >= 1 && v <= len(primes) {
				p := primes[v-1]
				if freq[p] > 0 {
					freq[p]--
					used = true
				}
			}
		}
		if !used {
			return fmt.Errorf("cannot validate value %d against remaining multiset", v)
		}
	}
	return nil
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
		// Validate candidate output by simulating removals from the multiset according to rules
		// Parse input into b list
		lines := strings.Split(strings.TrimSpace(input), "\n")
		if len(lines) < 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: malformed input fed to candidate\n", i)
			os.Exit(1)
		}
		// first line has n
		var nVal int
		fmt.Sscanf(strings.TrimSpace(lines[0]), "%d", &nVal)
		// second line has 2*n numbers
		parts := strings.Fields(lines[1])
		b := make([]int, 0, len(parts))
		for _, s := range parts {
			var v int
			fmt.Sscanf(s, "%d", &v)
			b = append(b, v)
		}
		if err := validateCandidate(got, nVal, b, primes, spfArr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed validation: %v\ninput:\n%s\nexpected example:%s\n   got:%s\n", i, err, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
