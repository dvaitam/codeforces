package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

const mod int64 = 1000000007

var fact []int64
var invFact []int64
var pow2 []int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func prep(maxN int) {
	fact = make([]int64, maxN+2)
	invFact = make([]int64, maxN+2)
	pow2 = make([]int64, maxN+2)
	fact[0] = 1
	pow2[0] = 1
	for i := 1; i <= maxN+1; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
		pow2[i] = pow2[i-1] * 2 % mod
	}
	invFact[maxN+1] = modPow(fact[maxN+1], mod-2)
	for i := maxN; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % mod
	}
}

func F2(n, m int) int64 {
	if n < 0 || m < 0 {
		return 0
	}
	limit := n
	if m < limit {
		limit = m
	}
	var res int64
	for r := 0; r <= limit; r++ {
		k := n + m - r
		term := fact[k+1]
		term = term * invFact[n-r] % mod
		term = term * invFact[m-r] % mod
		term = term * invFact[r] % mod
		term = term * pow2[k] % mod
		if r%2 == 1 {
			res -= term
		} else {
			res += term
		}
	}
	res %= mod
	if res < 0 {
		res += mod
	}
	return res
}

func solve(n, m int) int64 {
	ans := int64(0)
	ans = (ans - 3*F2(n-2, m-2)) % mod
	ans = (ans + 6*F2(n-2, m-1)) % mod
	ans = (ans - 3*F2(n-2, m)) % mod
	ans = (ans + 6*F2(n-1, m-2)) % mod
	ans = (ans - 8*F2(n-1, m-1)) % mod
	ans = (ans + 2*F2(n-1, m)) % mod
	ans = (ans - 3*F2(n, m-2)) % mod
	ans = (ans + 2*F2(n, m-1)) % mod
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func generateTests() (string, []int64) {
	rand.Seed(1)
	t := 100
	var buf bytes.Buffer
	tests := make([][2]int, t)
	fmt.Fprintln(&buf, t)
	maxSum := 0
	for i := 0; i < t; i++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(10) + 1
		tests[i] = [2]int{n, m}
		fmt.Fprintln(&buf, n, m)
		if n+m > maxSum {
			maxSum = n + m
		}
	}
	prep(maxSum)
	var answers []int64
	for _, p := range tests {
		answers = append(answers, solve(p[0], p[1]))
	}
	return buf.String(), answers
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	input, answers := generateTests()
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewBufferString(input)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for idx, exp := range answers {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(scanner.Text(), &got)
		if got != exp {
			fmt.Printf("case %d: expected %d, got %d\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
