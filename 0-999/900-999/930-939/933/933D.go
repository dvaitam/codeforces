package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const MOD int64 = 1000000007
const (
	inv2  int64 = 500000004
	inv6  int64 = 166666668
	inv30 int64 = 233333335
	inv42 int64 = 23809524
)

func sum2(n int64) int64 {
	n %= MOD
	return n * (n + 1) % MOD * (2*n + 1) % MOD * inv6 % MOD
}

func sum4(n int64) int64 {
	n0 := n % MOD
	a := n0 * (n0 + 1) % MOD * (2*n0 + 1) % MOD
	t := n0 * n0 % MOD
	b := (3*t%MOD + 3*n0 - 1) % MOD
	return a * b % MOD * inv30 % MOD
}

func sum6(n int64) int64 {
	n0 := n % MOD
	a := n0 * (n0 + 1) % MOD * (2*n0 + 1) % MOD
	t1 := n0 * n0 % MOD
	t2 := t1 * t1 % MOD
	t3 := t2 * 3 % MOD
	t4 := t1 * n0 % MOD
	t5 := t4 * 6 % MOD
	b := (t3 + t5 - 3*n0 + 1) % MOD
	return a * b % MOD * inv42 % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m int64
	if _, err := fmt.Fscan(in, &m); err != nil {
		return
	}

	sqrtM := int64(math.Sqrt(float64(m)))
	var S0, S1, S2, S3 int64
	for x := int64(0); x <= sqrtM; x++ {
		rX := x * x
		if rX > m {
			break
		}
		yMax := int64(math.Sqrt(float64(m - rX)))
		if x == 0 {
			S0 = (S0 + 1) % MOD
			continue
		}
		x2 := rX % MOD
		x4 := x2 * x2 % MOD
		x6 := x4 * x2 % MOD
		S0 = (S0 + 4) % MOD
		S1 = (S1 + 4*x2) % MOD
		S2 = (S2 + 4*x4) % MOD
		S3 = (S3 + 4*x6) % MOD
		Y := yMax
		if Y > x-1 {
			Y = x - 1
		}
		if Y >= 1 {
			cnt := Y % MOD
			sumY2 := sum2(Y)
			sumY4 := sum4(Y)
			sumY6 := sum6(Y)
			S0 = (S0 + 8*cnt%MOD) % MOD
			S1 = (S1 + 8*((cnt*x2%MOD+sumY2)%MOD)) % MOD
			S2 = (S2 + 8*((cnt*x4%MOD+2*x2*sumY2%MOD+sumY4)%MOD)) % MOD
			S3 = (S3 + 8*((cnt*x6%MOD+3*x4*sumY2%MOD+3*x2*sumY4%MOD+sumY6)%MOD)) % MOD
		}
		if yMax >= x {
			r := 2 * x * x
			rm := r % MOD
			S0 = (S0 + 4) % MOD
			S1 = (S1 + 4*rm%MOD) % MOD
			S2 = (S2 + 4*rm%MOD*rm%MOD) % MOD
			S3 = (S3 + 4*rm%MOD*rm%MOD*rm%MOD) % MOD
		}
	}
	mMod := m % MOD
	Tm := mMod * (mMod + 1) % MOD * inv2 % MOD
	Sm := mMod * (mMod + 1) % MOD * (2*mMod + 1) % MOD * inv6 % MOD
	part1 := (mMod + 1) % MOD * Tm % MOD
	B := (S2 - S1) % MOD
	if B < 0 {
		B += MOD
	}
	B = B * inv2 % MOD
	C := (2*S3%MOD - 3*S2%MOD + S1) % MOD
	if C < 0 {
		C += MOD
	}
	C = C * inv6 % MOD
	ans := (S0*(part1-Sm+MOD)%MOD - ((mMod+1)%MOD)*B%MOD + C) % MOD
	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(out, ans)
}
