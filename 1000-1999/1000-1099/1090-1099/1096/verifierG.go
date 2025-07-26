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

const modG = 998244353

func modAdd(a, b int) int {
	a += b
	if a >= modG {
		a -= modG
	}
	return a
}
func modSub(a, b int) int {
	a -= b
	if a < 0 {
		a += modG
	}
	return a
}
func modMul(a, b int) int { return int((int64(a) * int64(b)) % modG) }

func modPow(a, e int) int {
	res := 1
	x := a % modG
	for e > 0 {
		if e&1 != 0 {
			res = modMul(res, x)
		}
		x = modMul(x, x)
		e >>= 1
	}
	return res
}

func modInv(a int) int { return modPow(a, modG-2) }

func ntt(a []int, invert bool) {
	n := len(a)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for j&bit != 0 {
			j ^= bit
			bit >>= 1
		}
		j |= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		wlen := modPow(3, (modG-1)/length)
		if invert {
			wlen = modInv(wlen)
		}
		for i := 0; i < n; i += length {
			w := 1
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := modMul(a[i+j+half], w)
				a[i+j] = modAdd(u, v)
				a[i+j+half] = modSub(u, v)
				w = modMul(w, wlen)
			}
		}
	}
	if invert {
		invN := modInv(n)
		for i := range a {
			a[i] = modMul(a[i], invN)
		}
	}
}

func multiply(a, b []int) []int {
	need := len(a) + len(b) - 1
	n := 1
	for n < need {
		n <<= 1
	}
	fa := make([]int, n)
	fb := make([]int, n)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < n; i++ {
		fa[i] = modMul(fa[i], fb[i])
	}
	ntt(fa, true)
	return fa[:need]
}

func square(a []int) []int { return multiply(a, a) }

func solveG(n int, digits []int) int {
	allowed := [10]bool{}
	for _, d := range digits {
		allowed[d] = true
	}
	var ans func(int) []int
	ans = func(k int) []int {
		if k == 0 {
			return []int{1}
		}
		if k%2 == 0 {
			ret := ans(k / 2)
			return square(ret)
		}
		ret := ans(k - 1)
		res := make([]int, len(ret)+10)
		for i, v := range ret {
			if v == 0 {
				continue
			}
			for d := 0; d < 10; d++ {
				if allowed[d] {
					res[i+d] = modAdd(res[i+d], v)
				}
			}
		}
		return res
	}
	f := ans(n / 2)
	out := 0
	for _, v := range f {
		out = modAdd(out, modMul(v, v))
	}
	return out
}

type testCaseG struct {
	input    string
	expected int
}

func generateCaseG(rng *rand.Rand) testCaseG {
	n := 2 * (rng.Intn(4) + 1)
	k := rng.Intn(5) + 1
	digits := rng.Perm(10)[:k]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, d := range digits {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", d))
	}
	sb.WriteByte('\n')
	return testCaseG{input: sb.String(), expected: solveG(n, digits)}
}

func runCaseG(bin string, tc testCaseG) error {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCaseG{generateCaseG(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseG(rng))
	}
	for i, tc := range cases {
		if err := runCaseG(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
