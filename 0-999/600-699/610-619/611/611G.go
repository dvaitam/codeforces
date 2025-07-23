package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int64 = 1000000007

// Int128 represents a signed 128-bit integer using hi and lo parts.
type Int128 struct {
	hi int64
	lo uint64
}

func newInt128FromInt64(x int64) Int128 {
	return Int128{hi: x >> 63, lo: uint64(x)}
}

func (a Int128) Add(b Int128) Int128 {
	lo, c := bits.Add64(a.lo, b.lo, 0)
	hi := a.hi + b.hi + int64(c)
	return Int128{hi, lo}
}

func (a Int128) Sub(b Int128) Int128 {
	lo, brr := bits.Sub64(a.lo, b.lo, 0)
	hi := a.hi - b.hi - int64(brr)
	return Int128{hi, lo}
}

func (a Int128) Neg() Int128 {
	lo := ^a.lo + 1
	hi := ^a.hi
	if lo == 0 {
		hi++
	}
	return Int128{hi, lo}
}

func (a Int128) Cmp(b Int128) int {
	if a.hi == b.hi {
		if a.lo == b.lo {
			return 0
		}
		if a.lo < b.lo {
			return -1
		}
		return 1
	}
	if a.hi < b.hi {
		return -1
	}
	return 1
}

func (a Int128) Shl1() Int128 {
	hi := a.hi<<1 | int64(a.lo>>63)
	lo := a.lo << 1
	return Int128{hi, lo}
}

func (a Int128) MulInt64(b int64) Int128 {
	neg := false
	if a.hi < 0 {
		neg = !neg
		a = a.Neg()
	}
	if b < 0 {
		neg = !neg
		b = -b
	}
	p1hi, p1lo := bits.Mul64(a.lo, uint64(b))
	_, p2lo := bits.Mul64(uint64(a.hi), uint64(b))
	hi := p2lo + p1hi
	res := Int128{int64(hi), p1lo}
	if neg {
		res = res.Neg()
	}
	return res
}

func mul64(a, b int64) Int128 {
	neg := false
	if a < 0 {
		neg = !neg
		a = -a
	}
	if b < 0 {
		neg = !neg
		b = -b
	}
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	res := Int128{int64(hi), lo}
	if neg {
		res = res.Neg()
	}
	return res
}

var pow2_64_mod int64

func powMod(base, exp, mod int64) int64 {
	res := int64(1 % mod)
	b := base % mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * b % mod
		}
		b = b * b % mod
		exp >>= 1
	}
	return res
}

func (a Int128) Mod(mod int64) int64 {
	if mod == 0 {
		return 0
	}
	if a.hi < 0 {
		x := a.Neg().Mod(mod)
		if x == 0 {
			return 0
		}
		return mod - x
	}
	m := uint64(mod)
	hiPart := (uint64(a.hi) % m) * uint64(pow2_64_mod) % m
	loPart := a.lo % m
	return int64((hiPart + loPart) % m)
}

func cross(x1, y1, x2, y2 int64) int64 {
	return x1*y2 - x2*y1
}

func main() {
	pow2_64_mod = powMod(2, 64, MOD)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	X := make([]int64, 2*n+1)
	Y := make([]int64, 2*n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &X[i], &Y[i])
	}
	for i := n; i <= 2*n; i++ {
		X[i] = X[i-n]
		Y[i] = Y[i-n]
	}

	prefixX := make([]int64, 2*n+1)
	prefixY := make([]int64, 2*n+1)
	prefixCross := make([]Int128, 2*n+1)
	prefixCrossAcc := make([]Int128, 2*n+1)

	for i := 0; i < 2*n; i++ {
		prefixX[i+1] = prefixX[i] + X[i]
		prefixY[i+1] = prefixY[i] + Y[i]
		c := cross(X[i], Y[i], X[i+1], Y[i+1])
		prefixCross[i+1] = prefixCross[i].Add(newInt128FromInt64(c))
		prefixCrossAcc[i+1] = prefixCrossAcc[i].Add(prefixCross[i+1])
	}

	total := prefixCross[n]
	if total.hi < 0 {
		total = total.Neg()
	}
	totalMod := total.Mod(MOD)

	area2 := func(i, j int) Int128 {
		val := prefixCross[j].Sub(prefixCross[i])
		cp := cross(X[j], Y[j], X[i], Y[i])
		val = val.Add(newInt128FromInt64(cp))
		if val.hi < 0 {
			val = val.Neg()
		}
		return val
	}

	sumArea2 := func(i, l, r int) Int128 {
		if l > r {
			return Int128{}
		}
		S1 := prefixCrossAcc[r].Sub(prefixCrossAcc[l-1])
		len := int64(r - l + 1)
		S2 := prefixCross[i].MulInt64(len)
		sumX := prefixX[r+1] - prefixX[l]
		sumY := prefixY[r+1] - prefixY[l]
		t1 := mul64(Y[i], sumX)
		t2 := mul64(X[i], sumY)
		S3 := t1.Sub(t2)
		S := S1.Sub(S2).Add(S3)
		if S.hi < 0 {
			S = S.Neg()
		}
		return S
	}

	ans := int64(0)
	j := 0
	for i := 0; i < n; i++ {
		if j < i+1 {
			j = i + 1
		}
		end := i + n - 2
		for j+1 <= end {
			a2 := area2(i, j+1)
			if a2.Shl1().Cmp(total) <= 0 {
				j++
			} else {
				break
			}
		}
		start := i + 2
		if j < start-1 {
			j = start - 1
		}
		cnt1 := j - (start - 1)
		if cnt1 < 0 {
			cnt1 = 0
		}
		cnt2 := end - j
		if cnt2 < 0 {
			cnt2 = 0
		}
		sum1 := sumArea2(i, start, j)
		sum2 := sumArea2(i, j+1, end)
		ans = (ans + (totalMod*int64(cnt1))%MOD) % MOD
		ans = (ans + MOD - (2*sum1.Mod(MOD))%MOD) % MOD
		ans = (ans + (2*sum2.Mod(MOD))%MOD) % MOD
		ans = (ans + MOD - (totalMod*int64(cnt2))%MOD) % MOD
	}

	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(out, ans%MOD)
}
