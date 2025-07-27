package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solveCase(n int) string {
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
	mpf := make([]int, n+1)
	mpf[1] = 1
	for i := 2; i <= n; i++ {
		p := spf[i]
		mpf[i] = p
		if i/p > 1 {
			if mpf[i/p] > mpf[i] {
				mpf[i] = mpf[i/p]
			}
		}
	}
	freq := make([]int, n+1)
	for i := 1; i <= n; i++ {
		freq[mpf[i]]++
	}
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + freq[i]
	}
	primePre := make([]int, n+1)
	cnt := 0
	isPrime := make([]bool, n+1)
	for _, p := range primes {
		isPrime[p] = true
	}
	for i := 1; i <= n; i++ {
		if isPrime[i] {
			cnt++
		}
		primePre[i] = cnt
	}
	piN := cnt
	M := make([]int, n+1)
	for g := 1; g <= n; g++ {
		M[g] = prefix[g] + (piN - primePre[g])
	}
	res := make([]int, n+1)
	g := 1
	for k := 2; k <= n; k++ {
		for g <= n && M[g] < k {
			g++
		}
		if g > n {
			res[k] = n
		} else {
			res[k] = g
		}
	}
	var sb strings.Builder
	for k := 2; k <= n; k++ {
		if k > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(res[k]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildCase(n int) testCase {
	return testCase{input: fmt.Sprintf("%d\n", n), expected: solveCase(n)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 2
	return buildCase(n)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", tc.expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase(2),
		buildCase(3),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
