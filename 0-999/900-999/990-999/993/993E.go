package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

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
			w := complex(1, 0)
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
		invN := complex(1/float64(n), 0)
		for i := range a {
			a[i] *= invN
		}
	}
}

func convolutionInt64(a, b []int64) []int64 {
	n := 1
	needed := len(a) + len(b) - 1
	for n < needed {
		n <<= 1
	}
	fa := make([]complex128, n)
	fb := make([]complex128, n)
	for i := 0; i < len(a); i++ {
		fa[i] = complex(float64(a[i]), 0)
	}
	for i := 0; i < len(b); i++ {
		fb[i] = complex(float64(b[i]), 0)
	}
	fft(fa, false)
	fft(fb, false)
	for i := 0; i < n; i++ {
		fa[i] *= fb[i]
	}
	fft(fa, true)
	res := make([]int64, needed)
	for i := 0; i < needed; i++ {
		res[i] = int64(math.Round(real(fa[i])))
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var x int
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1]
		if arr[i-1] < x {
			prefix[i]++
		}
	}
	freq := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		freq[prefix[i]]++
	}
	rev := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		rev[i] = freq[n-i]
	}
	conv := convolutionInt64(freq, rev)
	ans := make([]int64, n+1)
	cross0 := conv[n]
	ans[0] = (cross0 - int64(n+1)) / 2
	for k := 1; k <= n; k++ {
		ans[k] = conv[n+k]
	}
	for i := 0; i <= n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(ans[i])
	}
	fmt.Println()
}
