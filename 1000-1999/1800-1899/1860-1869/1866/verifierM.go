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

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.ParseInt(expect, 10, 64)
	val, err := strconv.ParseInt(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output not integer: %v", err)
	}
	exp = ((exp % mod) + mod) % mod
	val = ((val % mod) + mod) % mod
	if exp != val {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([]int{0}),
		makeCase([]int{50, 50}),
		makeCase([]int{0, 0, 0}),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(5) + 1
		p := make([]int, n)
		for j := 0; j < n; j++ {
			p[j] = rand.Intn(100)
		}
		tests = append(tests, makeCase(p))
	}
	return tests
}

func makeCase(p []int) testCase {
	var sb strings.Builder
	n := len(p)
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(p)),
	}
}

func solveReference(p []int) int64 {
	inv100 := modPow(100, mod-2)
	pval := make([]int64, 100)
	for i := 0; i < 100; i++ {
		pval[i] = int64(i) * inv100 % mod
	}
	H := make([]int64, 100)
	s := int64(0)
	for _, val := range p {
		pr := pval[val]
		inv1 := modPow(int((1-pr+mod)%mod), mod-2)
		f := pr * H[val] % mod
		t := (inv1 + pr*inv1%mod*s%mod - f) % mod
		if t < 0 {
			t += mod
		}
		s = (s + t) % mod
		for r := 0; r < 100; r++ {
			H[r] = (H[r]*pval[r] + s) % mod
		}
	}
	return s % mod
}

func modPow(a, b int) int64 {
	x := int64(a)
	res := int64(1)
	y := int64(b)
	for y > 0 {
		if y&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		y >>= 1
	}
	return res
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
