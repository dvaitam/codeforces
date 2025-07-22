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

const MOD = 1000000009

func modPow(a, e int) int {
	res := 1
	x := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = int((int64(res) * int64(x)) % MOD)
		}
		x = int((int64(x) * int64(x)) % MOD)
		e >>= 1
	}
	return res
}

func computeA(M, h, k int) int {
	if k == 0 {
		return modPow(4, M)
	}
	if k == 4 {
		if M < h {
			return modPow(4, M)
		}
		return 0
	}
	dim := 1
	for i := 0; i < k; i++ {
		dim *= h
	}
	trans := make([][]int, k+1)
	for j := 0; j <= k; j++ {
		trans[j] = make([]int, dim)
	}
	decode := func(s int) []int {
		z := make([]int, k)
		for i := 0; i < k; i++ {
			z[i] = s % h
			s /= h
		}
		return z
	}
	encode := func(z []int) int {
		s := 0
		mul := 1
		for i := 0; i < k; i++ {
			s += z[i] * mul
			mul *= h
		}
		return s
	}
	for s := 0; s < dim; s++ {
		z := decode(s)
		for j := 0; j < k; j++ {
			ok := true
			nz := make([]int, k)
			for i := 0; i < k; i++ {
				if i == j {
					nz[i] = 0
				} else {
					if z[i]+1 >= h {
						ok = false
						break
					}
					nz[i] = z[i] + 1
				}
			}
			if ok {
				trans[j][s] = encode(nz)
			} else {
				trans[j][s] = -1
			}
		}
		ok := true
		nz := make([]int, k)
		for i := 0; i < k; i++ {
			if z[i]+1 >= h {
				ok = false
				break
			}
			nz[i] = z[i] + 1
		}
		if ok {
			trans[k][s] = encode(nz)
		} else {
			trans[k][s] = -1
		}
	}
	dp := make([]int, dim)
	ndp := make([]int, dim)
	dp[0] = 1
	others := 4 - k
	for pos := 0; pos < M; pos++ {
		for i := range ndp {
			ndp[i] = 0
		}
		for s := 0; s < dim; s++ {
			v := dp[s]
			if v == 0 {
				continue
			}
			for j := 0; j < k; j++ {
				t := trans[j][s]
				if t >= 0 {
					ndp[t] = (ndp[t] + v) % MOD
				}
			}
			t := trans[k][s]
			if t >= 0 {
				ndp[t] = (ndp[t] + int((int64(v)*int64(others))%MOD)) % MOD
			}
		}
		dp, ndp = ndp, dp
	}
	sum := 0
	for _, v := range dp {
		sum = (sum + v) % MOD
	}
	return sum
}

func solveCase(n, h int) int {
	M := n - h + 1
	A := make([]int, 5)
	for k := 0; k <= 4; k++ {
		A[k] = computeA(M, h, k)
	}
	C4 := []int{1, 4, 6, 4, 1}
	f := 0
	for k := 0; k <= 4; k++ {
		term := int((int64(C4[k]) * int64(A[k])) % MOD)
		if k%2 == 1 {
			f = (f - term + MOD) % MOD
		} else {
			f = (f + term) % MOD
		}
	}
	powTail := modPow(4, n-M)
	g := int((int64(f) * int64(powTail)) % MOD)
	total := modPow(4, n)
	ans := (total - g + MOD) % MOD
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	h := rng.Intn(minInt(n, 5)) + 1
	input := fmt.Sprintf("%d %d\n", n, h)
	expected := fmt.Sprintf("%d", solveCase(n, h))
	return input, expected
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
