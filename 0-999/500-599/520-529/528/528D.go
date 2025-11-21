package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	var sStr, tStr string
	fmt.Fscan(in, &sStr)
	fmt.Fscan(in, &tStr)

	s := []byte(sStr)
	t := []byte(tStr)

	good := make([][]int, 4)
	for c := 0; c < 4; c++ {
		diff := make([]int, n+1)
		for i := 0; i < n; i++ {
			if charIdx(s[i]) == c {
				l := i - k
				if l < 0 {
					l = 0
				}
				r := i + k
				if r >= n {
					r = n - 1
				}
				diff[l]++
				diff[r+1]--
			}
		}
		cur := 0
		good[c] = make([]int, n)
		for i := 0; i < n; i++ {
			cur += diff[i]
			if cur > 0 {
				good[c][i] = 1
			}
		}
	}

	totalLen := n + m - 1
	total := make([]int, totalLen)

	for c := 0; c < 4; c++ {
		vecS := make([]float64, n)
		for i := 0; i < n; i++ {
			vecS[i] = float64(good[c][i])
		}
		vecT := make([]float64, m)
		for j := 0; j < m; j++ {
			if charIdx(t[j]) == c {
				vecT[m-1-j] = 1
			}
		}
		conv := convolution(vecS, vecT)
		for i := 0; i < totalLen; i++ {
			total[i] += conv[i]
		}
	}

	ans := 0
	for start := 0; start <= n-m; start++ {
		if total[start+m-1] == m {
			ans++
		}
	}

	fmt.Fprintln(out, ans)
}

func charIdx(b byte) int {
	switch b {
	case 'A':
		return 0
	case 'T':
		return 1
	case 'G':
		return 2
	case 'C':
		return 3
	default:
		return -1
	}
}

func convolution(a, b []float64) []int {
	if len(a) == 0 || len(b) == 0 {
		return []int{}
	}
	need := len(a) + len(b) - 1
	n := 1
	for n < need {
		n <<= 1
	}

	fa := make([]complex128, n)
	fb := make([]complex128, n)
	for i := 0; i < len(a); i++ {
		fa[i] = complex(a[i], 0)
	}
	for i := 0; i < len(b); i++ {
		fb[i] = complex(b[i], 0)
	}

	fft(fa, false)
	fft(fb, false)

	for i := 0; i < n; i++ {
		fa[i] *= fb[i]
	}

	fft(fa, true)

	res := make([]int, need)
	for i := 0; i < need; i++ {
		res[i] = int(math.Round(real(fa[i])))
	}
	return res
}

func fft(a []complex128, invert bool) {
	n := len(a)
	for i, j := 1, 0; i < n; i++ {
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
		angle := 2 * math.Pi / float64(length)
		if invert {
			angle = -angle
		}
		wlen := complex(math.Cos(angle), math.Sin(angle))
		for i := 0; i < n; i += length {
			w := complex(1, 0)
			for j := 0; j < length/2; j++ {
				u := a[i+j]
				v := a[i+j+length/2] * w
				a[i+j] = u + v
				a[i+j+length/2] = u - v
				w *= wlen
			}
		}
	}

	if invert {
		ni := complex(float64(1)/float64(n), 0)
		for i := 0; i < n; i++ {
			a[i] *= ni
		}
	}
}
