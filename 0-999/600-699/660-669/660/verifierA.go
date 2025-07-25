package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveA(n int, arr []int) (int, []int) {
	const maxP = 1000
	isPrime := make([]bool, maxP+1)
	for i := 2; i <= maxP; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= maxP; i++ {
		if isPrime[i] {
			for j := i * i; j <= maxP; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= maxP; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	b := make([]int, n)
	ans := 0
	for i := 0; i+1 < n; i++ {
		if gcd(arr[i], arr[i+1]) != 1 {
			for _, p := range primes {
				if gcd(arr[i], p) == 1 && gcd(p, arr[i+1]) == 1 {
					b[i] = p
					break
				}
			}
			ans++
		}
	}
	res := make([]int, 0, n+ans)
	for i := 0; i < n; i++ {
		res = append(res, arr[i])
		if i < n-1 && b[i] != 0 {
			res = append(res, b[i])
		}
	}
	return ans, res
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(8) + 2
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rand.Intn(30) + 1
		}
		var sb strings.Builder
		fmt.Fprintln(&sb, n)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		ans, res := solveA(n, append([]int(nil), arr...))
		var exp strings.Builder
		fmt.Fprintln(&exp, ans)
		for i, v := range res {
			if i > 0 {
				exp.WriteByte(' ')
			}
			fmt.Fprint(&exp, v)
		}
		exp.WriteByte('\n')
		output, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(output)
		want := strings.TrimSpace(exp.String())
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", t+1, sb.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
