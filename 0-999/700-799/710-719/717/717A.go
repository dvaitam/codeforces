package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000007

// ext represents a + b*sqrt(5) in F_p[sqrt(5)]
type ext struct{ a, b int64 }

func (x ext) add(y ext) ext { return ext{(x.a + y.a) % MOD, (x.b + y.b) % MOD} }
func (x ext) sub(y ext) ext { return ext{(x.a - y.a + MOD) % MOD, (x.b - y.b + MOD) % MOD} }
func (x ext) mul(y ext) ext {
	// (a+b√5)*(c+d√5) = (ac + 5bd) + (ad+bc)√5
	return ext{
		(x.a*y.a + 5*x.b%MOD*y.b) % MOD,
		(x.a*y.b + x.b*y.a) % MOD,
	}
}
func (x ext) neg() ext { return ext{(MOD - x.a) % MOD, (MOD - x.b) % MOD} }

// inv computes multiplicative inverse in ext
func (x ext) inv() ext {
	// 1/(a+b√5) = (a - b√5)/(a^2 - 5b^2)
	d := (x.a*x.a - 5*(x.b*x.b)%MOD + MOD) % MOD
	invd := modinv(d)
	return ext{(x.a * invd) % MOD, ((MOD - x.b) * invd) % MOD}
}
func (x ext) pow(e int64) ext {
	var res ext = ext{1, 0}
	base := x
	for e > 0 {
		if e&1 == 1 {
			res = res.mul(base)
		}
		base = base.mul(base)
		e >>= 1
	}
	return res
}

func modinv(x int64) int64 {
	return powmod(x, MOD-2)
}
func powmod(a, e int64) int64 {
	var r int64 = 1
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			r = (r * a) % MOD
		}
		a = (a * a) % MOD
		e >>= 1
	}
	return r
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	var l, r int64
	fmt.Fscan(in, &k, &l, &r)
	// precompute f(n) for small n and find threshold n0 where f(n)>=k
	smallF := []int64{0, 2, 3} // f(1)=2,f(2)=3
	n0 := int64(1)
	for {
		if n0 < 2 {
			if smallF[n0] >= int64(k) {
				break
			}
		} else {
			x := smallF[n0] + smallF[n0-1]
			smallF = append(smallF, x)
			if x >= int64(k) {
				break
			}
		}
		n0++
	}
	// sum small n naive
	var ans int64
	up := r
	if up >= n0 {
		up = n0 - 1
	}
	for n := l; n <= up; n++ {
		// compute f(n)
		fn := smallF[n]
		if fn >= int64(k) {
			ans = (ans + C(fn, int64(k))) % MOD
		}
	}
	if r < n0 {
		fmt.Println(ans)
		return
	}
	L := max64(l, n0)
	// prepare polynomial p(x) = prod_{j=0..k-1}(x - j) = sum_{d=0..k} p[d]*x^d
	p := make([]int64, k+1)
	p[0] = 1
	for j := 0; j < k; j++ {
		for d := j + 1; d >= 1; d-- {
			p[d] = (p[d-1] - int64(j)*p[d]%MOD + MOD) % MOD
		}
	}
	invfk := modinv(fact(int64(k)))
	// P coefficients
	P := make([]int64, k+1)
	for d := 0; d <= k; d++ {
		P[d] = p[d] * invfk % MOD
	}
	// precompute binomials C[n][i]
	Cb := make([][]int64, k+1)
	for n := 0; n <= k; n++ {
		Cb[n] = make([]int64, n+1)
		Cb[n][0], Cb[n][n] = 1, 1
		for i := 1; i < n; i++ {
			Cb[n][i] = (Cb[n-1][i-1] + Cb[n-1][i]) % MOD
		}
	}
	// compute alpha, beta, A, C
	inv2 := modinv(2)
	inv5 := modinv(5)
	sqrt5 := ext{0, 1}
	_ = sqrt5
	inv_sqrt5 := ext{0, inv5}
	alpha := ext{inv2, inv2}
	beta := ext{inv2, (MOD - inv2) % MOD}
	a2 := alpha.mul(alpha)
	b2 := beta.mul(beta)
	A := a2.mul(inv_sqrt5)
	Cc := b2.mul(inv_sqrt5).neg()
	// precompute powers
	powA := make([]ext, k+1)
	powC := make([]ext, k+1)
	powA[0], powC[0] = ext{1, 0}, ext{1, 0}
	for i := 1; i <= k; i++ {
		powA[i] = powA[i-1].mul(A)
		powC[i] = powC[i-1].mul(Cc)
	}
	// powers of alpha and beta
	powAl := make([]ext, k+1)
	powBe := make([]ext, k+1)
	powAl[0], powBe[0] = ext{1, 0}, ext{1, 0}
	for i := 1; i <= k; i++ {
		powAl[i] = powAl[i-1].mul(alpha)
		powBe[i] = powBe[i-1].mul(beta)
	}
	// sum over i,j
	var sum ext
	total := r - L + 1
	for i := 0; i <= k; i++ {
		for j := 0; j+i <= k; j++ {
			coeff := P[i+j] * Cb[i+j][i] % MOD
			termBase := powA[i].mul(powC[j])
			coeffE := ext{coeff, 0}.mul(termBase)
			d := powAl[i].mul(powBe[j])
			// geometric sum of d^n for n=L..r = d^L * (1 - d^total) / (1 - d)
			dL := d.pow(L)
			num := ext{1, 0}.sub(d.pow(total))
			denom := ext{1, 0}.sub(d)
			geom := num.mul(denom.inv()).mul(dL)
			sum = sum.add(coeffE.mul(geom))
		}
	}
	// sum.a is result for n>=L
	ans = (ans + sum.a) % MOD
	fmt.Println(ans)
}

func fact(n int64) int64 {
	var r int64 = 1
	for i := int64(1); i <= n; i++ {
		r = r * i % MOD
	}
	return r
}

func C(n, k int64) int64 {
	if n < k || k < 0 {
		return 0
	}
	var r int64 = 1
	for i := int64(0); i < k; i++ {
		r = r * ((n - i) % MOD) % MOD
	}
	r = r * modinv(fact(k)) % MOD
	return r
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
