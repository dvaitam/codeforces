package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod = 1000000007

type test struct{ a, b, n int }

func modPow(x, p int64) int64 {
	res := int64(1)
	x %= mod
	for p > 0 {
		if p&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		p >>= 1
	}
	return res
}

func countExcellent(a, b, n int) int64 {
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	var res int64
	for k := 0; k <= n; k++ {
		sum := k*a + (n-k)*b
		x := sum
		good := true
		for x > 0 {
			d := x % 10
			if d != a && d != b {
				good = false
				break
			}
			x /= 10
		}
		if !good {
			continue
		}
		comb := fact[n] * invFact[k] % mod * invFact[n-k] % mod
		res = (res + comb) % mod
	}
	return res
}

func genTests() []test {
	rand.Seed(3)
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		a := rand.Intn(8) + 1
		b := rand.Intn(9-a) + a + 1
		n := rand.Intn(30) + 1
		tests = append(tests, test{a, b, n})
	}
	return tests
}

func buildInput(t test) string {
	return fmt.Sprintf("%d %d %d\n", t.a, t.b, t.n)
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func verifyOutput(out string, t test) bool {
	out = strings.TrimSpace(out)
	val, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return false
	}
	expected := countExcellent(t.a, t.b, t.n)
	return val%mod == expected%mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(cand, input)
		if err != nil {
			fmt.Printf("test %d: run error %v\n", i+1, err)
			os.Exit(1)
		}
		if !verifyOutput(out, t) {
			fmt.Printf("test %d failed. input:%soutput:%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
