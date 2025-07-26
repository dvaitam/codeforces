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

const mod = 1000000007

func solve(primes []int) string {
	seen := make(map[int]bool)
	prod := int64(1)
	for _, p := range primes {
		if !seen[p] {
			seen[p] = true
			prod = (prod * int64(p)) % mod
		}
	}
	return fmt.Sprintf("%d\n", prod)
}

func buildCase(pr []int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(pr)))
	sb.WriteByte('\n')
	for i, v := range pr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solve(pr)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	pr := make([]int, n)
	for i := range pr {
		pr[i] = randPrime(rng)
	}
	return buildCase(pr)
}

func randPrime(rng *rand.Rand) int {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}
	return primes[rng.Intn(len(primes))]
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase([]int{2, 3, 5}),
		buildCase([]int{7, 7, 7, 11}),
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
