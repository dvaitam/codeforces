package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// fft performs in-place Fast Fourier Transform on a.
func fft(a []complex128, invert bool) {
	n := len(a)
	// bit-reversal permutation
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
		invN := complex(1/float64(n), 0)
		for i := 0; i < n; i++ {
			a[i] *= invN
		}
	}
}

func nextPow2(n int) int {
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}

// DSU for 6 elements
type dsu struct{ p [6]int }

func newDSU() *dsu {
	d := &dsu{}
	for i := 0; i < 6; i++ {
		d.p[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) unite(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	d.p[rb] = ra
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var S, T string
	if _, err := fmt.Fscan(reader, &S, &T); err != nil {
		return
	}
	n := len(S)
	m := len(T)
	if m == 0 {
		for i := 0; i <= n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, 0)
		}
		writer.WriteByte('\n')
		return
	}
	limit := nextPow2(n + m)

	// prepare FFT arrays for S characters
	fs := make([][]complex128, 6)
	for ch := 0; ch < 6; ch++ {
		arr := make([]complex128, limit)
		b := byte('a' + ch)
		for i := 0; i < n; i++ {
			if S[i] == b {
				arr[i] = 1
			}
		}
		fft(arr, false)
		fs[ch] = arr
	}
	// prepare FFT arrays for reversed T characters
	ft := make([][]complex128, 6)
	for ch := 0; ch < 6; ch++ {
		arr := make([]complex128, limit)
		b := byte('a' + ch)
		for i := 0; i < m; i++ {
			if T[m-1-i] == b {
				arr[i] = 1
			}
		}
		fft(arr, false)
		ft[ch] = arr
	}

	// mapping for edge indices
	pair := make([][2]int, 0, 15)
	idx := [6][6]int{}
	cnt := 0
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {
			idx[i][j] = cnt
			idx[j][i] = cnt
			pair = append(pair, [2]int{i, j})
			cnt++
		}
	}
	total := n - m + 1
	edges := make([][]bool, len(pair))
	for i := range edges {
		edges[i] = make([]bool, total)
	}

	tmp := make([]complex128, limit)
	for a := 0; a < 6; a++ {
		for b := a + 1; b < 6; b++ {
			for i := 0; i < limit; i++ {
				tmp[i] = fs[a][i]*ft[b][i] + fs[b][i]*ft[a][i]
			}
			fft(tmp, true)
			arr := edges[idx[a][b]]
			base := m - 1
			for i := 0; i < total; i++ {
				if real(tmp[i+base]) > 0.5 {
					arr[i] = true
				}
			}
			for i := 0; i < limit; i++ {
				tmp[i] = 0
			}
		}
	}

	res := make([]int, total)
	for pos := 0; pos < total; pos++ {
		d := newDSU()
		ops := 0
		for idxEdge, pr := range pair {
			if edges[idxEdge][pos] {
				if d.unite(pr[0], pr[1]) {
					ops++
				}
			}
		}
		res[pos] = ops
	}
	for i := 0; i < total; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, res[i])
	}
	writer.WriteByte('\n')
}
