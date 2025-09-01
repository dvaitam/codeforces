package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(cmdPath string, input string) (string, error) {
	cmd := exec.Command(cmdPath)
	if strings.HasSuffix(cmdPath, ".go") {
		cmd = exec.Command("go", "run", cmdPath)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// semantic oracle for 1670F
const mod1670 = int64(1000000007)

func modPow1670(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod1670
		}
		a = a * a % mod1670
		e >>= 1
	}
	return res
}

func countUpTo(n int, x, z int64, choose []int64) int64 {
	if x < 0 {
		return 0
	}
	maxBit := 61
	dp := make([][2]int64, n+1)
	dp[0][0] = 1
	for i := 0; i < maxBit; i++ {
		xi := int((x >> i) & 1)
		zi := int((z >> i) & 1)
		next := make([][2]int64, n+1)
		for c := 0; c <= n; c++ {
			for less := 0; less < 2; less++ {
				ways := dp[c][less]
				if ways == 0 {
					continue
				}
				bitSum := (c & 1) ^ zi
				if less == 0 && bitSum > xi {
					continue
				}
				start := (c - bitSum + 1) / 2
				if start < 0 {
					start = 0
				}
				end := (c + n - bitSum) / 2
				if end > n {
					end = n
				}
				newLessBase := less
				if less == 0 && bitSum < xi {
					newLessBase = 1
				}
				for cp := start; cp <= end; cp++ {
					k := 2*cp + bitSum - c
					if k < 0 || k > n {
						continue
					}
					val := ways * choose[k] % mod1670
					next[cp][newLessBase] = (next[cp][newLessBase] + val) % mod1670
				}
			}
		}
		dp = next
	}
	return (dp[0][0] + dp[0][1]) % mod1670
}

func expected1670F(n int, l, r, z int64) int64 {
	// precompute C(n, k)
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod1670
	}
	invFact[n] = modPow1670(fact[n], mod1670-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod1670
	}
	choose := make([]int64, n+1)
	for k := 0; k <= n; k++ {
		choose[k] = fact[n] * invFact[k] % mod1670 * invFact[n-k] % mod1670
	}
	ans := countUpTo(n, r, z, choose) - countUpTo(n, l-1, z, choose)
	ans %= mod1670
	if ans < 0 {
		ans += mod1670
	}
	return ans
}

func genTests() []string {
	rand.Seed(6)
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		l := rand.Int63n(50) + 1
		r := l + rand.Int63n(50)
		z := rand.Int63n(64)
		// Single test format: n l r z\n
		tests = append(tests, fmt.Sprintf("%d %d %d %d\n", n, l, r, z))
	}
	// Edge cases
	tests = append(tests, "1 1 1 0\n")
	tests = append(tests, "2 1 3 1\n")
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		// compute expected
		fields := strings.Fields(tc)
		if len(fields) < 4 {
			fmt.Printf("bad test %d: %q\n", i+1, tc)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		l, _ := strconv.ParseInt(fields[1], 10, 64)
		r, _ := strconv.ParseInt(fields[2], 10, 64)
		z, _ := strconv.ParseInt(fields[3], 10, 64)
		expVal := expected1670F(n, l, r, z)
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr = strings.TrimSpace(gotStr)
		gotFields := strings.Fields(gotStr)
		if len(gotFields) == 0 {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%d\ngot:%q\n", i+1, tc, expVal, gotStr)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(gotFields[0], 10, 64)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%d\ngot:%q\n", i+1, tc, expVal, gotStr)
			os.Exit(1)
		}
		if gotVal != expVal {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%d\ngot:%d\n", i+1, tc, expVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
