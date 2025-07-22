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

const mod = 1000000007

func isLucky(x int) bool {
	if x <= 0 {
		return false
	}
	for x > 0 {
		d := x % 10
		if d != 4 && d != 7 {
			return false
		}
		x /= 10
	}
	return true
}

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow(a, mod-2) }

func nCr(n int64, r int, fact, invFact []int64) int64 {
	if r < 0 || n < int64(r) {
		return 0
	}
	return fact[n] * invFact[r] % mod * invFact[n-int64(r)] % mod
}

func solveCase(n, k int, arr []int) int64 {
	freq := make(map[int]int)
	U := 0
	for _, x := range arr {
		if isLucky(x) {
			freq[x]++
		} else {
			U++
		}
	}
	counts := make([]int, 0, len(freq))
	for _, c := range freq {
		counts = append(counts, c)
	}
	m := len(counts)
	dp := make([]int64, m+1)
	dp[0] = 1
	for _, c := range counts {
		cc := int64(c)
		for t := m; t >= 1; t-- {
			dp[t] = (dp[t] + dp[t-1]*cc) % mod
		}
	}
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modInv(fact[n])
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	ans := int64(0)
	maxT := m
	if k < maxT {
		maxT = k
	}
	for t := 0; t <= maxT; t++ {
		rem := k - t
		if rem < 0 || rem > U {
			continue
		}
		waysLucky := dp[t]
		waysUnlucky := nCr(int64(U), rem, fact, invFact)
		ans = (ans + waysLucky*waysUnlucky) % mod
	}
	return ans
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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Intn(4) == 0 {
			// lucky number
			digits := []int{4, 7}
			val := 0
			for j := 0; j < rng.Intn(3)+1; j++ {
				val = val*10 + digits[rng.Intn(2)]
			}
			arr[i] = val
		} else {
			arr[i] = rng.Intn(100) + 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", solveCase(n, k, arr))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
