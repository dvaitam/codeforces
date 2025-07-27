package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 998244353

func modexp(a, e int) int {
	res := 1
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		e >>= 1
	}
	return res
}

func modinv(a int) int {
	return modexp((a%mod+mod)%mod, mod-2)
}

func solveCase(n, m int) int {
	fact := make([]int, n+1)
	invfact := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
	}
	invfact[n] = modinv(fact[n])
	for i := n; i > 0; i-- {
		invfact[i-1] = int(int64(invfact[i]) * int64(i) % mod)
	}
	powA := make([]int, n+1)
	powB := make([]int, n+1)
	A := (m + 1) % mod
	B := m % mod
	powA[0], powB[0] = 1, 1
	for i := 1; i <= n; i++ {
		powA[i] = int(int64(powA[i-1]) * int64(A) % mod)
		powB[i] = int(int64(powB[i-1]) * int64(B) % mod)
	}
	S := 0
	for j := 1; j <= n; j++ {
		bin := int(int64(fact[n]) * int64(invfact[j]) % mod * int64(invfact[n-j]) % mod)
		aj := powA[j]
		bj := powB[j]
		denom := aj - bj
		if denom < 0 {
			denom += mod
		}
		invd := modinv(denom)
		term := int(int64(bin) * int64(aj) % mod * int64(invd) % mod)
		if j%2 == 1 {
			S += term
			if S >= mod {
				S -= mod
			}
		} else {
			S -= term
			if S < 0 {
				S += mod
			}
		}
	}
	coeff := int(int64(n+m+1) % mod * int64(modinv(m+1)) % mod)
	ans := int(int64(S) * int64(coeff) % mod)
	return ans
}

func generateTest() (string, string) {
	n := rand.Intn(5) + 1
	m := rand.Intn(5) + 1
	inp := fmt.Sprintf("%d %d\n", n, m)
	out := fmt.Sprintf("%d\n", solveCase(n, m))
	return inp, out
}

func referenceIO(t int) (string, string) {
	var in strings.Builder
	var out strings.Builder
	for i := 0; i < t; i++ {
		ti, to := generateTest()
		in.WriteString(ti)
		out.WriteString(to)
	}
	return in.String(), out.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	rand.Seed(8)
	in, exp := referenceIO(100)
	out, err := runBinary(os.Args[1], in)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Println("Wrong Answer")
		fmt.Println("Expected:\n" + exp)
		fmt.Println("Got:\n" + out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
