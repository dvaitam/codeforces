package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt64() int64 {
	sign := int64(1)
	var val int64
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func egcd(a, b int64) (int64, int64, int64) {
	if b == 0 {
		if a >= 0 {
			return a, 1, 0
		}
		return -a, -1, 0
	}
	g, x1, y1 := egcd(b, a%b)
	x := y1
	y := x1 - (a/b)*y1
	return g, x, y
}

func modInverse(a, mod int64) int64 {
	if mod == 1 {
		return 0
	}
	g, x, _ := egcd(a, mod)
	if g != 1 && g != -1 {
		return 0
	}
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func solveLinear(total, pref, mod int64) (bool, int64, int64) {
	if mod == 0 {
		if pref == 0 {
			return true, 0, 1
		}
		return false, 0, 0
	}
	c := ((-pref)%mod + mod) % mod
	if total == 0 {
		if c == 0 {
			return true, 0, 1
		}
		return false, 0, 0
	}
	g := gcd(total, mod)
	if c%g != 0 {
		return false, 0, 0
	}
	modOut := mod / g
	if modOut == 1 {
		return true, 0, 1
	}
	a := (total / g) % modOut
	if a < 0 {
		a += modOut
	}
	cDiv := c / g
	inv := modInverse(a, modOut)
	rem := (inv * cDiv) % modOut
	rem = (rem%modOut + modOut) % modOut
	return true, rem, modOut
}

func crt(r1, m1, r2, m2 int64) (bool, int64, int64) {
	if m1 == 1 {
		return true, r2 % m2, m2
	}
	if m2 == 1 {
		return true, r1 % m1, m1
	}
	g := gcd(m1, m2)
	diff := r2 - r1
	if diff%g != 0 {
		return false, 0, 0
	}
	lcm := m1 / g * m2
	if lcm == 0 {
		return false, 0, 0
	}
	m2g := m2 / g
	k1 := m1 / g
	diffDiv := diff / g
	diffDiv %= m2g
	if diffDiv < 0 {
		diffDiv += m2g
	}
	inv := modInverse(k1%m2g, m2g)
	x := (diffDiv * inv) % m2g
	r := r1 + m1*x
	r %= lcm
	if r < 0 {
		r += lcm
	}
	return true, r, lcm
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt64()
	for ; t > 0; t-- {
		n := int(in.NextInt64())
		k := in.NextInt64()
		w := in.NextInt64()
		h := in.NextInt64()
		var s string
		fmt.Fscan(in.r, &s)
		modX := 2 * w
		modY := 2 * h
		prefX := make([]int64, n+1)
		prefY := make([]int64, n+1)
		for i := 0; i < n; i++ {
			ch := s[i]
			switch ch {
			case 'L':
				prefX[i+1] = prefX[i] - 1
				prefY[i+1] = prefY[i]
			case 'R':
				prefX[i+1] = prefX[i] + 1
				prefY[i+1] = prefY[i]
			case 'U':
				prefX[i+1] = prefX[i]
				prefY[i+1] = prefY[i] + 1
			case 'D':
				prefX[i+1] = prefX[i]
				prefY[i+1] = prefY[i] - 1
			}
		}
		sumX := prefX[n]
		sumY := prefY[n]
		var ans int64
		for i := 1; i <= n; i++ {
			ok1, r1, m1 := solveLinear(sumX, prefX[i], modX)
			if !ok1 {
				continue
			}
			ok2, r2, m2 := solveLinear(sumY, prefY[i], modY)
			if !ok2 {
				continue
			}
			sol, r, lcm := crt(r1, m1, r2, m2)
			if !sol {
				continue
			}
			first := r % lcm
			if first < 0 {
				first += lcm
			}
			if first > k-1 {
				continue
			}
			cnt := (k-1-first)/lcm + 1
			ans += cnt
		}
		fmt.Fprintln(out, ans)
	}
}
