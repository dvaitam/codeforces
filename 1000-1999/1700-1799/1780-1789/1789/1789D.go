package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		var n int
		var sa, sb string
		fmt.Fscan(reader, &n, &sa, &sb)
		origA := []byte(sa)
		bArr := []byte(sb)
		var v []int

		// define recursive functions
		var left, right func(a []byte, bit int) bool

		right = func(a []byte, bit int) bool {
			if bytes.Equal(a, bArr) {
				return true
			}
			n := len(a)
			i := n - 1
			for i >= 0 && a[i] != '1' {
				i--
			}
			if i < 0 {
				if bit == 0 {
					return false
				}
				// switch to left with bit=0
				// copy a to avoid modifying for caller
				a2 := make([]byte, n)
				copy(a2, a)
				return left(a2, 0)
			}
			for j := i - 1; j >= 0; j-- {
				if a[j] != bArr[j] {
					k := i - j
					v = append(v, k)
					// build x: a shifted right by k
					x := make([]byte, n)
					for r := 0; r < n; r++ {
						if r+k < n {
							x[r] = a[r+k]
						} else {
							x[r] = '0'
						}
					}
					// a = a XOR x
					for r := 0; r < n; r++ {
						ar := a[r] - '0'
						xr := x[r] - '0'
						a[r] = '0' + ((ar ^ xr) & 1)
					}
				}
			}
			if bit == 1 {
				return left(a, 0)
			}
			return bytes.Equal(a, bArr)
		}

		left = func(a []byte, bit int) bool {
			if bytes.Equal(a, bArr) {
				return true
			}
			n := len(a)
			i := 0
			for i < n && a[i] != '1' {
				i++
			}
			if i >= n {
				if bit == 1 {
					a2 := make([]byte, n)
					copy(a2, a)
					return right(a2, 0)
				}
				return false
			}
			for j := i + 1; j < n; j++ {
				if a[j] != bArr[j] {
					k := j - i
					v = append(v, -k)
					// build x: a shifted left by k
					x := make([]byte, n)
					for r := 0; r < n; r++ {
						if r < k {
							x[r] = '0'
						} else {
							x[r] = a[r-k]
						}
					}
					// a = a XOR x
					for r := 0; r < n; r++ {
						ar := a[r] - '0'
						xr := x[r] - '0'
						a[r] = '0' + ((ar ^ xr) & 1)
					}
				}
			}
			if bit == 1 {
				return right(a, 0)
			}
			return bytes.Equal(a, bArr)
		}

		// attempt left then right
		a1 := left(func() []byte { tmp := make([]byte, n); copy(tmp, origA); return tmp }(), 1)
		if a1 {
			fmt.Fprintln(writer, len(v))
			if len(v) > 0 {
				for _, x := range v {
					fmt.Fprint(writer, x, " ")
				}
				fmt.Fprintln(writer)
			}
		} else {
			// reset v and try right
			v = v[:0]
			a2 := make([]byte, n)
			copy(a2, origA)
			a1 = right(a2, 1)
			if a1 {
				fmt.Fprintln(writer, len(v))
				if len(v) > 0 {
					for _, x := range v {
						fmt.Fprint(writer, x, " ")
					}
					fmt.Fprintln(writer)
				}
			} else {
				fmt.Fprintln(writer, -1)
			}
		}
	}
}
