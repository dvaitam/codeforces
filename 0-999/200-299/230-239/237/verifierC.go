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

type testCase struct {
	input    string
	expected int
}

func isPrimeSieve(limit int) []bool {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	return isPrime
}

func prefixPrimes(isPrime []bool) []int {
	pre := make([]int, len(isPrime))
	for i := 1; i < len(isPrime); i++ {
		pre[i] = pre[i-1]
		if isPrime[i] {
			pre[i]++
		}
	}
	return pre
}

func ok(a, b, k, l int, pre []int) bool {
	for x := a; x <= b-l+1; x++ {
		cnt := pre[x+l-1] - pre[x-1]
		if cnt < k {
			return false
		}
	}
	return true
}

func solve(a, b, k int) int {
	n := b - a + 1
	isPrime := isPrimeSieve(b)
	pre := prefixPrimes(isPrime)
	lo, hi := k, n
	ans := -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if ok(a, b, k, mid, pre) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) testCase {
	a := rng.Intn(1000) + 1
	b := a + rng.Intn(1000)
	k := rng.Intn(b-a+1) + 1
	l := solve(a, b, k)
	input := fmt.Sprintf("%d %d %d\n", a, b, k)
	return testCase{input: input, expected: l}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
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
	var cases []testCase
	cases = append(cases, testCase{input: "1 1 1\n", expected: 1})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
