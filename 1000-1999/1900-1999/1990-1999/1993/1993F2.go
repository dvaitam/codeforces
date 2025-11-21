package main

import (
	"bufio"
	"fmt"
	"os"
)

type axisParam struct {
	mod    int64 // A = 2*w or 2*h
	total  int64
	g      int64
	free   bool
	redMod int64
	inv    int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func extGcd(a, b int64) (int64, int64, int64) {
	if b == 0 {
		if a < 0 {
			return -a, -1, 0
		}
		return a, 1, 0
	}
	g, x1, y1 := extGcd(b, a%b)
	x, y := y1, x1-(a/b)*y1
	return g, x, y
}

func modInverse(a, mod int64) int64 {
	if mod == 1 {
		return 0
	}
	g, x, _ := extGcd(a, mod)
	if g != 1 && g != -1 {
		return 0
	}
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func normalize(x, mod int64) int64 {
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func prepareAxis(total, mod int64) axisParam {
	total = normalize(total, mod)
	free := total == 0
	if free {
		return axisParam{mod: mod, total: total, g: mod, free: true}
	}
	g := gcd(total, mod)
	redMod := mod / g
	inv := modInverse(total/g%redMod, redMod)
	return axisParam{mod: mod, total: total, g: g, redMod: redMod, inv: inv}
}

func axisConstraint(pref int64, param axisParam) (bool, bool, int64, int64) {
	pref = normalize(pref, param.mod)
	if param.free {
		if pref != 0 {
			return false, false, 0, 1
		}
		return true, true, 0, 1
	}
	if pref%param.g != 0 {
		return false, false, 0, 1
	}
	redPref := pref / param.g
	rhs := normalize(-redPref, param.redMod)
	rem := (rhs * param.inv) % param.redMod
	return true, false, rem, param.redMod
}

func countSolutions(r, mod, k int64) int64 {
	if mod == 1 {
		if r >= k {
			return 0
		}
		return k
	}
	r = normalize(r, mod)
	if r >= k {
		return 0
	}
	return 1 + (k-1-r)/mod
}

func crt(r1, m1, r2, m2 int64) (bool, int64, int64) {
	g := gcd(m1, m2)
	if (r2-r1)%g != 0 {
		return false, 0, 0
	}
	lcm := m1 / g * m2
	if m2/g == 1 {
		r := normalize(r1, lcm)
		return true, r, lcm
	}
	m1Div := m1 / g
	m2Div := m2 / g
	rhs := (r2 - r1) / g
	rhs = normalize(rhs, m2Div)
	inv := modInverse(m1Div%m2Div, m2Div)
	s := (rhs * inv) % m2Div
	res := r1 + m1*s
	res = normalize(res, lcm)
	return true, res, lcm
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		var w, h int64
		fmt.Fscan(in, &n, &k, &w, &h)
		var s string
		fmt.Fscan(in, &s)

		prefixX := make([]int64, n+1)
		prefixY := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			prefixX[i] = prefixX[i-1]
			prefixY[i] = prefixY[i-1]
			switch s[i-1] {
			case 'L':
				prefixX[i]--
			case 'R':
				prefixX[i]++
			case 'D':
				prefixY[i]--
			case 'U':
				prefixY[i]++
			}
		}

		ax := prepareAxis(prefixX[n], 2*w)
		ay := prepareAxis(prefixY[n], 2*h)

		var ans int64
		for i := 1; i <= n; i++ {
			okx, freex, remx, modx := axisConstraint(prefixX[i], ax)
			if !okx {
				continue
			}
			oky, freey, remy, mody := axisConstraint(prefixY[i], ay)
			if !oky {
				continue
			}
			if freex && freey {
				ans += k
			} else if freex {
				ans += countSolutions(remy, mody, k)
			} else if freey {
				ans += countSolutions(remx, modx, k)
			} else {
				ok, rem, mod := crt(remx, modx, remy, mody)
				if !ok {
					continue
				}
				ans += countSolutions(rem, mod, k)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
