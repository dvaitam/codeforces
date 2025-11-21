package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int64 {
	sign := int64(1)
	val := int64(0)
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

var bigOne = big.NewInt(1)

func computeQ(a []int64) (int64, bool) {
	n := len(a)
	C := big.NewInt(0)
	tmp := big.NewInt(0)
	for i := n - 1; i >= 0; i-- {
		C.Lsh(C, 1)
		tmp.SetInt64(a[i])
		C.Add(C, tmp)
	}

	M := big.NewInt(1)
	M.Lsh(M, uint(n))
	M.Sub(M, bigOne)

	rem := new(big.Int).Mod(C, M)
	if rem.Sign() != 0 {
		return 0, false
	}

	qBig := new(big.Int).Div(C, M)
	return qBig.Int64(), true
}

func feasible(a []int64, q, x int64) bool {
	y := q - x
	if y < 0 {
		return false
	}
	carry := y
	for _, val := range a {
		cur := val + carry - x
		if cur < 0 || (cur&1) == 1 {
			return false
		}
		carry = cur / 2
	}
	return carry == y
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := int(in.NextInt())
		a := make([]int64, n)
		sum := int64(0)
		for i := 0; i < n; i++ {
			a[i] = in.NextInt()
			sum += a[i]
		}

		q, ok := computeQ(a)
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}

		avg := sum / int64(n)
		upper := avg
		if q < upper {
			upper = q
		}

		lo, hi := int64(0), upper
		for lo < hi {
			mid := (lo + hi + 1) >> 1
			if feasible(a, q, mid) {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		ans := sum - int64(n)*lo
		fmt.Fprintln(out, ans)
	}
}
