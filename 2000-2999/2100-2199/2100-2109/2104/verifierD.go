package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const maxN = 400000
const sieveLimit = 7000000

var primePrefix [maxN + 1]int64

func init() {
	primes := make([]int, 0, maxN)
	isComposite := make([]bool, sieveLimit+1)
	for i := 2; i <= sieveLimit && len(primes) < maxN; i++ {
		if !isComposite[i] {
			primes = append(primes, i)
			if i*i <= sieveLimit {
				for j := i * i; j <= sieveLimit; j += i {
					isComposite[j] = true
				}
			}
		}
	}
	if len(primes) < maxN {
		panic("not enough primes in sieve")
	}
	for i := 1; i <= maxN; i++ {
		primePrefix[i] = primePrefix[i-1] + int64(primes[i-1])
	}
}

func expectedRemovals(arr []int64) int {
	n := len(arr)
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + arr[i-1]
	}
	best := 0
	for l := n; l >= 0; l-- {
		if primePrefix[l] <= prefix[l] {
			best = l
			break
		}
	}
	return n - best
}

type testCase struct {
	arr []int64
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(1000) + 2)
	}
	return testCase{arr: arr}
}

func buildInput(cases []testCase) (string, []int) {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	exp := make([]int, len(cases))
	for i, tc := range cases {
		fmt.Fprintln(&sb, len(tc.arr))
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte("\n")
		exp[i] = expectedRemovals(tc.arr)
	}
	return sb.String(), exp
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/2104D_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input, expected := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) < len(expected) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, exp := range expected {
		got, err := strconv.Atoi(lines[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on case %d: %q\n", i+1, lines[i])
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %d\narray: %v\n", i+1, exp, got, cases[i].arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
