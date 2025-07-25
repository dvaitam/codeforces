package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const maxValA = 1000000

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func primesUpTo(n int) []int {
	sieve := make([]bool, n+1)
	var primes []int
	for i := 2; i <= n; i++ {
		if !sieve[i] {
			primes = append(primes, i)
			for j := i * 2; j <= n; j += i {
				sieve[j] = true
			}
		}
	}
	return primes
}

var primeListA = primesUpTo(1000000)

func findPrimeStep(a []int) int {
	diff := make([]bool, maxValA+1)
	n := len(a)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			d := a[i] - a[j]
			if d < 0 {
				d = -d
			}
			if d <= maxValA {
				diff[d] = true
			}
		}
	}
	for _, p := range primeListA {
		ok := true
		for x := p; x <= maxValA; x += p {
			if diff[x] {
				ok = false
				break
			}
		}
		if ok {
			return p
		}
	}
	return 0
}

func checkCaseA(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var t, n int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return err
	}
	if t != 1 {
		return errors.New("expected one test case")
	}
	if _, err := fmt.Fscan(in, &n); err != nil {
		return err
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(in, &a[i]); err != nil {
			return err
		}
	}
	expectPossible := findPrimeStep(a) > 0

	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return errors.New("no output")
	}
	if tokens[0] == "NO" {
		if expectPossible {
			return errors.New("expected YES but got NO")
		}
		if len(tokens) != 1 {
			return errors.New("extra data after NO")
		}
		return nil
	}
	if tokens[0] != "YES" {
		return errors.New("first token should be YES or NO")
	}
	if !expectPossible {
		return errors.New("expected NO but got YES")
	}
	if len(tokens) != 1+n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens)-1)
	}
	b := make([]int, n)
	used := map[int]bool{}
	for i := 0; i < n; i++ {
		x, err := strconv.Atoi(tokens[1+i])
		if err != nil {
			return fmt.Errorf("bad number %q", tokens[1+i])
		}
		if x < 1 || x > maxValA {
			return fmt.Errorf("number %d out of range", x)
		}
		if used[x] {
			return fmt.Errorf("duplicate number %d", x)
		}
		used[x] = true
		b[i] = x
	}
	sums := map[int]bool{}
	for _, x := range a {
		for _, y := range b {
			s := x + y
			if sums[s] {
				return fmt.Errorf("duplicate sum %d", s)
			}
			sums[s] = true
		}
	}
	return nil
}

func generateCaseA(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := rng.Intn(maxValA) + 1
		sb.WriteString(strconv.Itoa(val))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []string
	// some fixed cases
	cases = append(cases, "1\n1\n1\n")
	cases = append(cases, "1\n2\n1 2\n")
	cases = append(cases, "1\n3\n1 2 3\n")

	for len(cases) < 100 {
		cases = append(cases, generateCaseA(rng))
	}

	for i, tc := range cases {
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if err := checkCaseA(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
