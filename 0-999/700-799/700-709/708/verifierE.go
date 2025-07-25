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

const MOD = 1000000007

func modAdd(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}
func modSub(a, b int) int {
	a -= b
	if a < 0 {
		a += MOD
	}
	return a
}
func modMul(a, b int) int { return int(int64(a) * int64(b) % MOD) }
func modPow(a, e int) int {
	res := 1
	x := a
	for e > 0 {
		if e&1 != 0 {
			res = modMul(res, x)
		}
		x = modMul(x, x)
		e >>= 1
	}
	return res
}
func modInv(a int) int { return modPow(a, MOD-2) }

func solveCase(n, m, a, b, k int) int {
	invB := modInv(b)
	p := modMul(a, invB)
	q := modSub(1, p)
	F := make([]int, m)
	F[0] = modPow(q, k)
	invQ := modInv(q)
	for i := 1; i < m; i++ {
		num := (k - (i - 1)) % MOD
		F[i] = modMul(F[i-1], num)
		F[i] = modMul(F[i], modInv(i))
		F[i] = modMul(F[i], p)
		F[i] = modMul(F[i], invQ)
	}
	PF := make([]int, m)
	PF[0] = F[0]
	for i := 1; i < m; i++ {
		PF[i] = modAdd(PF[i-1], F[i])
	}
	S0 := modSub(1, PF[m-1])
	size := m + 1
	dp1 := make([]int, size*size)
	dp2 := make([]int, size*size)
	for l := 0; l <= m; l++ {
		Pl := 0
		if l < m {
			Pl = F[l]
		} else {
			Pl = S0
		}
		if Pl == 0 {
			continue
		}
		maxR := m - 1 - l
		for r := 0; r <= maxR; r++ {
			Pr := 0
			if r < m-l {
				Pr = F[r]
			} else {
				if m-l-1 >= 0 {
					Pr = modSub(1, PF[m-l-1])
				} else {
					Pr = 1
				}
			}
			dp1[l*size+r] = modMul(Pl, Pr)
		}
	}
	for row := 2; row <= n; row++ {
		for i := 0; i < size; i++ {
			sum := 0
			for j := 0; j < size; j++ {
				sum = modAdd(sum, dp1[i*size+j])
				dp2[i*size+j] = sum
			}
		}
		for j := 0; j < size; j++ {
			for i := 1; i < size; i++ {
				dp2[i*size+j] = modAdd(dp2[i*size+j], dp2[(i-1)*size+j])
			}
		}
		for idx := range dp1 {
			dp1[idx] = 0
		}
		for l := 0; l <= m; l++ {
			Pl := 0
			if l < m {
				Pl = F[l]
			} else {
				Pl = S0
			}
			if Pl == 0 {
				continue
			}
			for r := 0; r <= m-l-1; r++ {
				Pr := 0
				if r < m-l {
					Pr = F[r]
				} else {
					if m-l-1 >= 0 {
						Pr = modSub(1, PF[m-l-1])
					} else {
						Pr = 1
					}
				}
				x := m - 1 - r
				y := m - 1 - l
				if x >= m {
					x = m
				}
				if y >= m {
					y = m
				}
				sum := dp2[x*size+y]
				if sum != 0 {
					dp1[l*size+r] = modAdd(dp1[l*size+r], modMul(modMul(Pl, Pr), sum))
				}
			}
		}
	}
	ans := 0
	for _, v := range dp1 {
		ans = modAdd(ans, v)
	}
	return ans
}

func runCase(bin string, n, m, a, b, k int) error {
	input := fmt.Sprintf("%d %d\n%d %d\n%d\n", n, m, a, b, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := fmt.Sprintf("%d", solveCase(n, m, a, b, k))
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) (int, int, int, int, int) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	a := rng.Intn(9) + 1
	b := rng.Intn(9) + a // ensure a<=b
	k := rng.Intn(5)
	return n, m, a, b, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, m, a, b, k := randomCase(rng)
		if err := runCase(bin, n, m, a, b, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
