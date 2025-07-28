package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n   int
	arr []int
}

func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func factorize(x int, primes []int, counts map[int]int) {
	for _, p := range primes {
		if p*p > x {
			break
		}
		for x%p == 0 {
			counts[p]++
			x /= p
		}
	}
	if x > 1 {
		counts[x]++
	}
}

func expected(tc Test, primes []int) int {
	counts := make(map[int]int)
	for _, v := range tc.arr {
		factorize(v, primes, counts)
	}
	pairs := 0
	leftover := 0
	for _, c := range counts {
		pairs += c / 2
		leftover += c % 2
	}
	return pairs + leftover/3
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(3)
	const cases = 100
	tests := make([]Test, cases)
	for i := range tests {
		n := rand.Intn(4) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(50) + 2
		}
		tests[i] = Test{n: n, arr: arr}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d\n", tc.n)
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("error running binary:", err)
		fmt.Print(out.String())
		return
	}

	primes := sieve(100)
	reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
	for idx, tc := range tests {
		var ans int
		if _, err := fmt.Fscan(reader, &ans); err != nil {
			fmt.Printf("test %d: failed to read output\n", idx+1)
			return
		}
		exp := expected(tc, primes)
		if ans != exp {
			fmt.Printf("test %d: expected %d got %d\n", idx+1, exp, ans)
			return
		}
	}
	fmt.Printf("verified %d test cases\n", len(tests))
}
