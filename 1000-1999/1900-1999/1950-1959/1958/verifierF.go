package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD int = 1_000_000_007

func modPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = int(int64(res) * int64(a) % int64(MOD))
		}
		a = int(int64(a) * int64(a) % int64(MOD))
		b >>= 1
	}
	return res
}

func combPrecompute(maxN int) ([]int, []int) {
	fac := make([]int, maxN+1)
	ifac := make([]int, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = int(int64(fac[i-1]) * int64(i) % int64(MOD))
	}
	ifac[maxN] = modPow(fac[maxN], MOD-2)
	for i := maxN; i > 0; i-- {
		ifac[i-1] = int(int64(ifac[i]) * int64(i) % int64(MOD))
	}
	return fac, ifac
}

func C(n, r int, fac, ifac []int) int {
	if r < 0 || r > n {
		return 0
	}
	return int(int64(fac[n]) * int64(ifac[r]) % int64(MOD) * int64(ifac[n-r]) % int64(MOD))
}

func solve(n, k int) []int {
	maxF := 2 * n
	fac, ifac := combPrecompute(maxF)

	ans := make([]int, n+1)
	for s := 0; s <= 2*n-2; s++ {
		val := 0
		if s <= 2*n-4 && k >= 2 && k-2 <= 2*n-4-s {
			a := 0
			if s-(n-2) > a {
				a = s - (n - 2)
			}
			b := s
			if n-2 < b {
				b = n - 2
			}
			if b >= a {
				cnt := b - a + 1
				v := C(2*n-4-s, k-2, fac, ifac)
				val = (val + int(int64(cnt)*int64(v)%int64(MOD))) % MOD
			}
		}
		if s >= n-1 && s <= 2*n-3 && k >= 1 && k-1 <= 2*n-3-s {
			v := C(2*n-3-s, k-1, fac, ifac)
			val = (val + int((2*int64(v))%int64(MOD))) % MOD
		}
		if s < n-2 {
			ans[0] = (ans[0] + val) % MOD
		} else {
			idx := s - (n - 2)
			if idx <= n {
				ans[idx] = (ans[idx] + val) % MOD
			}
		}
	}
	for i := 0; i <= n; i++ {
		ans[i] %= MOD
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(47)
	// Use smaller n for speed
	n := rand.Intn(20) + 2
	k := rand.Intn(2*n-2) + 2

	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, k)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input.String())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}

	gotParts := strings.Fields(strings.TrimSpace(string(outBytes)))
	if len(gotParts) != n+1 {
		fmt.Printf("expected %d numbers, got %d\n", n+1, len(gotParts))
		os.Exit(1)
	}
	want := solve(n, k)
	for i, s := range gotParts {
		var g int
		fmt.Sscan(s, &g)
		if g%MOD != want[i]%MOD {
			fmt.Printf("mismatch at index %d expected %d got %d\n", i, want[i]%MOD, g%MOD)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
