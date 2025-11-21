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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &t)

	sBytes := []byte(s)
	tBytes := []byte(t)

	mLen := n - m + 1
	if mLen < 0 {
		fmt.Fprintln(out, 0)
		fmt.Fprintln(out)
		return
	}

	arr1 := make([]float64, n)
	arr2 := make([]float64, n)
	for i := 0; i < n; i++ {
		val := float64(sBytes[i] - 'a' + 1)
		arr1[i] = val
		arr2[i] = val * val
	}

	cnt := make([]int, 26)
	pos := make([][]int, 26)
	for i := 0; i < m; i++ {
		idx := int(tBytes[i] - 'a')
		cnt[idx]++
		pos[idx] = append(pos[idx], i)
	}

	length := 1
	for length < n+m-1 {
		length <<= 1
	}

	base := make([]complex128, length)
	for i := 0; i < n; i++ {
		base[i] = complex(arr1[i], arr2[i])
	}
	fft(base, false)

	assign := make([][]int, 26)
	for c := 0; c < 26; c++ {
		if cnt[c] == 0 {
			continue
		}
		assign[c] = make([]int, mLen)
		buf := make([]complex128, length)
		for _, idx := range pos[c] {
			buf[m-1-idx] = complex(1, 0)
		}
		fft(buf, false)
		for i := 0; i < length; i++ {
			buf[i] *= base[i]
		}
		fft(buf, true)
		for i := 0; i < mLen; i++ {
			idx := i + m - 1
			sum1 := int(math.Round(real(buf[idx])))
			sum2 := int(math.Round(imag(buf[idx])))
			if sum1%cnt[c] != 0 {
				assign[c][i] = -1
				continue
			}
			val := sum1 / cnt[c]
			if val < 1 || val > 26 {
				assign[c][i] = -1
				continue
			}
			if sum2 != cnt[c]*val*val {
				assign[c][i] = -1
				continue
			}
			assign[c][i] = val - 1
		}
	}

	mapTo := make([]int, 26)
	targetOwner := make([]int, 26)
	ans := make([]int, 0)

	for i := 0; i < mLen; i++ {
		for k := 0; k < 26; k++ {
			mapTo[k] = -1
			targetOwner[k] = -1
		}
		ok := true
		for c := 0; c < 26 && ok; c++ {
			if cnt[c] == 0 {
				continue
			}
			target := assign[c][i]
			if target == -1 {
				ok = false
				break
			}
			mapTo[c] = target
			if targetOwner[target] == -1 {
				targetOwner[target] = c
			} else if targetOwner[target] != c {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		for c := 0; c < 26 && ok; c++ {
			if cnt[c] == 0 {
				continue
			}
			target := mapTo[c]
			if target >= 0 && cnt[target] > 0 {
				if mapTo[target] != c {
					ok = false
					break
				}
			}
		}
		if ok {
			ans = append(ans, i+1)
		}
	}

	fmt.Fprintln(out, len(ans))
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
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
