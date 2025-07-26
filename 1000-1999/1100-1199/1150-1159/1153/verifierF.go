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

const mod = 998244353

type Test struct {
	input    string
	expected string
}

func add(a, b int) int {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
}
func sub(a, b int) int {
	a -= b
	if a < 0 {
		a += mod
	}
	return a
}
func mul(a, b int) int { return int((int64(a) * int64(b)) % mod) }
func powmod(a, e int) int {
	res := 1
	x := a
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, x)
		}
		x = mul(x, x)
		e >>= 1
	}
	return res
}

func solve(n, k int, l int64) int {
	maxN := 2*n + 1
	fact := make([]int, maxN+1)
	invf := make([]int, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = mul(fact[i-1], i)
	}
	invf[maxN] = powmod(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invf[i-1] = mul(invf[i], i)
	}
	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = add(pow2[i-1], pow2[i-1])
	}
	cnr := make([]int, n+1)
	for r := 0; r <= n; r++ {
		cnr[r] = mul(fact[n], mul(invf[r], invf[n-r]))
	}
	Fr := make([]int, n+1)
	for r := 0; r <= n; r++ {
		Fr[r] = mul(pow2[r], mul(fact[r], mul(fact[r], invf[2*r+1])))
	}
	sumS := 0
	for i := 0; i < k; i++ {
		for r := i; r <= n; r++ {
			coef := mul(cnr[r], mul(fact[r], mul(invf[i], invf[r-i])))
			term := mul(coef, Fr[r])
			if (r-i)&1 == 1 {
				term = sub(0, term)
			}
			sumS = add(sumS, term)
		}
	}
	res := int((l % int64(mod)) * int64(sub(1, sumS)) % int64(mod))
	if res < 0 {
		res += mod
	}
	return res
}

func generateTests() []Test {
	rng := rand.New(rand.NewSource(47))
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		k := rng.Intn(n) + 1
		l := rng.Int63n(1000000) + 1
		input := fmt.Sprintf("%d %d %d\n", n, k, l)
		expected := solve(n, k, l)
		tests = append(tests, Test{input: input, expected: strconv.Itoa(expected)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%s\ngot:%s\n", i+1, tc.input, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
