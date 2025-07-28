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

func runCandidate(bin, input string) (string, error) {
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

const mod = 998244353

func modAdd(a, b int) int {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
}
func modSub(a, b int) int {
	a -= b
	if a < 0 {
		a += mod
	}
	return a
}
func modMul(a, b int) int { return int(int64(a) * int64(b) % int64(mod)) }
func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 != 0 {
			res = modMul(res, a)
		}
		a = modMul(a, a)
		e >>= 1
	}
	return res
}
func modInv(a int) int { return modPow(a, mod-2) }

func ntt(a []int, invert bool) {
	n := len(a)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j |= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		wlen := modPow(3, (mod-1)/length)
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

func solveCase(n, k, f int) int {
	if f > 2*k {
		return 0
	}
	g := make([][]int, n+1)
	g[1] = make([]int, k+1)
	for i := 0; i <= k; i++ {
		g[1][i] = 1
	}
	for depth := 2; depth <= n; depth++ {
		s := multiply(g[depth-1], g[depth-1])
		big := 0
		for i := k + 1; i < len(s); i++ {
			big = modAdd(big, s[i])
		}
		prefix := make([]int, k+2)
		sum := 0
		for i := k; i >= 0; i-- {
			sum = modAdd(sum, s[i])
			prefix[i] = sum
		}
		arr := make([]int, k+1)
		for x := 0; x <= k; x++ {
			val := big
			if x+1 <= k {
				val = modAdd(val, prefix[x+1])
			}
			contrib := modMul(s[x], (k-x+1)%mod)
			val = modAdd(val, contrib)
			arr[x] = val
		}
		g[depth] = arr
	}
	rootConv := multiply(g[n], g[n])
	if f >= len(rootConv) {
		return 0
	}
	return rootConv[f] % mod
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(4) + 1 // 1..4
	k := rng.Intn(8) + 1 // 1..8
	f := rng.Intn(2*k + 1)
	input := fmt.Sprintf("%d %d %d\n", n, k, f)
	expect := solveCase(n, k, f)
	return input, expect
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != fmt.Sprintf("%d", expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, expect, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
