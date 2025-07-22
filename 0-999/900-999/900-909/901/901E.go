package main

import (
	"bufio"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
)

// nextPow2 returns the smallest power of two >= n
func nextPow2(n int) int {
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}

// fft performs in-place radix-2 FFT or inverse FFT depending on invert flag.
func fft(a []complex128, invert bool) {
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
		angle := 2 * math.Pi / float64(length)
		if invert {
			angle = -angle
		}
		wlen := complex(math.Cos(angle), math.Sin(angle))
		for i := 0; i < n; i += length {
			w := complex(1.0, 0.0)
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := a[i+j+half] * w
				a[i+j] = u + v
				a[i+j+half] = u - v
				w *= wlen
			}
		}
	}
	if invert {
		invN := complex(1.0/float64(n), 0)
		for i := 0; i < n; i++ {
			a[i] *= invN
		}
	}
}

// convolution computes convolution of a and b using FFT
func convolution(a, b []complex128) []complex128 {
	n := len(a) + len(b) - 1
	size := nextPow2(n)
	fa := make([]complex128, size)
	fb := make([]complex128, size)
	copy(fa, a)
	copy(fb, b)
	fft(fa, false)
	fft(fb, false)
	for i := 0; i < size; i++ {
		fa[i] *= fb[i]
	}
	fft(fa, true)
	return fa[:n]
}

// bluestein computes the DFT of a using Bluestein's algorithm.
// If invert is true, computes the inverse transform.
func bluestein(a []complex128, invert bool) []complex128 {
	n := len(a)
	if n == 1 {
		res := make([]complex128, 1)
		res[0] = a[0]
		return res
	}
	m := nextPow2(2*n - 1)
	A := make([]complex128, m)
	B := make([]complex128, m)
	sign := -1.0
	if invert {
		sign = 1.0
	}
	factor := math.Pi / float64(n)
	for i := 0; i < n; i++ {
		angle := sign * float64(i*i) * factor
		A[i] = a[i] * cmplx.Exp(complex(0, angle))
		B[i] = cmplx.Exp(complex(0, -sign*float64(i*i)*factor))
		if i != 0 {
			B[m-i] = B[i]
		}
	}
	C := convolution(A, B)
	res := make([]complex128, n)
	for i := 0; i < n; i++ {
		angle := sign * float64(i*i) * factor
		res[i] = C[i+n-1] * cmplx.Exp(complex(0, angle))
		if invert {
			res[i] /= complex(float64(n), 0)
		}
	}
	return res
}

// solveLinear solves for x in M*x=b where M is circulant defined by bVec.
func solveLinear(bVec []float64, rhs []float64) []float64 {
	n := len(bVec)
	a := make([]complex128, n)
	for i := 0; i < n; i++ {
		a[i] = complex(rhs[i], 0)
	}
	b := make([]complex128, n)
	for i := 0; i < n; i++ {
		b[i] = complex(bVec[i], 0)
	}
	Bfft := bluestein(b, false)
	for i := 0; i < n; i++ {
		if cmplx.Abs(Bfft[i]) == 0 {
			// Should not happen due to independence
			return nil
		}
	}
	Afft := bluestein(a, false)
	for i := 0; i < n; i++ {
		Afft[i] /= cmplx.Conj(Bfft[i])
	}
	resC := bluestein(Afft, true)
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = real(resC[i])
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	c := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &c[i])
	}
	for i := 1; i < n; i++ {
		if (c[i]-c[0])%2 != 0 {
			fmt.Println(0)
			return
		}
	}
	d := make([]float64, n)
	for i := 0; i < n; i++ {
		d[i] = float64(c[0]-c[i]) / 2
	}
	baseA := solveLinear(b, d)
	if baseA == nil {
		fmt.Println(0)
		return
	}
	e := make([]float64, n)
	for i := range e {
		e[i] = 1
	}
	baseV := solveLinear(b, e)
	if baseV == nil {
		fmt.Println(0)
		return
	}
	var sv2, sa2, sav float64
	for i := 0; i < n; i++ {
		sv2 += baseV[i] * baseV[i]
		sa2 += baseA[i] * baseA[i]
		sav += baseA[i] * baseV[i]
	}
	Bsum2 := 0
	for i := 0; i < n; i++ {
		Bsum2 += int(b[i] * b[i])
	}
	Acoef := sv2
	Bcoef := 2 * (sav - 1)
	Ccoef := sa2 - float64(c[0]-Bsum2)
	disc := Bcoef*Bcoef - 4*Acoef*Ccoef
	if disc < -1e-6 {
		fmt.Println(0)
		return
	}
	if disc < 0 {
		disc = 0
	}
	sqrtD := math.Sqrt(disc)
	rs := []float64{(-Bcoef + sqrtD) / (2 * Acoef)}
	if sqrtD > 1e-8 {
		rs = append(rs, (-Bcoef-sqrtD)/(2*Acoef))
	}
	solutions := [][]int{}
	for _, r := range rs {
		a := make([]int, n)
		S := 0.0
		ok := true
		for i := 0; i < n; i++ {
			val := baseA[i] + r*baseV[i]
			iv := math.Round(val)
			if math.Abs(val-iv) > 1e-5 {
				ok = false
				break
			}
			a[i] = int(iv)
			S += iv * iv
		}
		if !ok {
			continue
		}
		// verify
		for i := 0; i < n; i++ {
			Bi := r + d[i]
			cc := float64(Bsum2) - 2*Bi + S
			if math.Abs(cc-float64(c[i])) > 1e-3 {
				ok = false
				break
			}
		}
		if ok {
			solutions = append(solutions, a)
		}
	}
	if len(solutions) == 0 {
		fmt.Println(0)
		return
	}
	// deduplicate and sort
	uniq := map[string]struct{}{}
	uniqSol := [][]int{}
	for _, sol := range solutions {
		key := fmt.Sprint(sol)
		if _, exist := uniq[key]; !exist {
			uniq[key] = struct{}{}
			uniqSol = append(uniqSol, sol)
		}
	}
	sort.Slice(uniqSol, func(i, j int) bool {
		a := uniqSol[i]
		b := uniqSol[j]
		for k := 0; k < n; k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return false
	})
	fmt.Println(len(uniqSol))
	for _, sol := range uniqSol {
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(sol[i])
		}
		fmt.Println()
	}
}
