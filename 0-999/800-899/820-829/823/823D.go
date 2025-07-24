package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func fft(a []complex128, invert bool) {
	n := len(a)
	j := 0
	for i := 1; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		ang := 2 * math.Pi / float64(length)
		if invert {
			ang = -ang
		}
		wlen := complex(math.Cos(ang), math.Sin(ang))
		for i := 0; i < n; i += length {
			w := complex(1, 0)
			half := length / 2
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
		invN := complex(1/float64(n), 0)
		for i := range a {
			a[i] *= invN
		}
	}
}

func convolution(a, b []complex128) []complex128 {
	n := 1
	for n < len(a)+len(b) {
		n <<= 1
	}
	fa := make([]complex128, n)
	fb := make([]complex128, n)
	copy(fa, a)
	copy(fb, b)
	fft(fa, false)
	fft(fb, false)
	for i := 0; i < n; i++ {
		fa[i] *= fb[i]
	}
	fft(fa, true)
	return fa
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		a := make([]complex128, n)
		b := make([]complex128, n)
		c := make([]complex128, n)
		d := make([]complex128, n)
		for i := 0; i < n; i++ {
			if s[i] == 'V' {
				a[i] = 1
				d[n-1-i] = 1
			}
			if s[i] == 'K' {
				b[n-1-i] = 1
				c[i] = 1
			}
		}
		convVK := convolution(a, b)
		convKV := convolution(c, d)
		bad := make([]bool, n)
		for diff := 1; diff < n; diff++ {
			idx := n - 1 - diff
			var x, y float64
			if idx >= 0 && idx < len(convVK) {
				x = real(convVK[idx])
				y = real(convKV[idx])
			}
			if int(math.Round(x)) > 0 || int(math.Round(y)) > 0 {
				bad[diff] = true
			}
		}
		valid := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			valid[i] = true
		}
		for i := 1; i <= n; i++ {
			for j := i; j < n; j += i {
				if bad[j] {
					valid[i] = false
					break
				}
			}
		}
		ans := []int{}
		for i := 1; i <= n; i++ {
			if valid[i] {
				ans = append(ans, i)
			}
		}
		fmt.Fprintln(writer, len(ans))
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
