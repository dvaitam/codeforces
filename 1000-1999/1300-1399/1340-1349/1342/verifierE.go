package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 998244353

func powmod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func initFact(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = powmod(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	return fact, invFact
}

func C(fact, inv []int64, n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * inv[r] % mod * inv[n-r] % mod
}

func surjection(fact, inv []int64, n, c int) int64 {
	res := int64(0)
	for i := 0; i <= c; i++ {
		term := C(fact, inv, c, i) * powmod(int64(c-i), int64(n)) % mod
		if i%2 == 1 {
			res = (res - term) % mod
		} else {
			res = (res + term) % mod
		}
	}
	if res < 0 {
		res += mod
	}
	return res
}

type testCase struct {
	n, k int
}

func (tc testCase) Input() string {
	return fmt.Sprintf("%d %d\n", tc.n, tc.k)
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solve(tc testCase) string {
	n, k := tc.n, tc.k
	if k > n-1 {
		return "0"
	}
	fact, inv := initFact(n)
	if k == 0 {
		return fmt.Sprint(fact[n] % mod)
	}
	c := n - k
	val := C(fact, inv, n, c) * surjection(fact, inv, n, c) % mod
	ans := val * 2 % mod
	return fmt.Sprint(ans)
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(5))
	tests := make([]testCase, 100)
	for i := range tests {
		n := rng.Intn(20) + 1
		k := rng.Intn(n + 1)
		tests[i] = testCase{n: n, k: k}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want := solve(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if strings.TrimSpace(out) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
