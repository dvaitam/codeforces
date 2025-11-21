package main

import (
	"bufio"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b == '\n') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if b == '-' {
		sign = -1
		b, err = fs.r.ReadByte()
	}
	for err == nil && b >= '0' && b <= '9' {
		val = val*10 + int(b-'0')
		b, err = fs.r.ReadByte()
	}
	return sign * val
}

func (fs *fastScanner) nextInt64() int64 {
	sign, val := int64(1), int64(0)
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b == '\n') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if b == '-' {
		sign = -1
		b, err = fs.r.ReadByte()
	}
	for err == nil && b >= '0' && b <= '9' {
		val = val*10 + int64(b-'0')
		b, err = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		l := fs.nextInt64()
		_ = fs.nextInt64() // r, equals l in this version

		a := make([]int, n+1)
		pre := make([]int, n+1)
		for i := 1; i <= n; i++ {
			a[i] = fs.nextInt()
			pre[i] = pre[i-1] ^ a[i]
		}

		nEvenTerm := 0
		if n%2 == 0 {
			nEvenTerm = pre[n/2]
		}

		memo := make(map[int64]int)
		var pref func(int64) int
		pref = func(k int64) int {
			if k <= int64(n) {
				return pre[int(k)]
			}
			if val, ok := memo[k]; ok {
				return val
			}
			ans := pre[n]
			if n%2 == 0 {
				ans ^= nEvenTerm
			}
			if k%2 == 0 {
				ans ^= pref(k / 2)
			}
			memo[k] = ans
			return ans
		}

		var get func(int64) int
		get = func(idx int64) int {
			if idx <= int64(n) {
				return a[int(idx)]
			}
			return pref(idx / 2)
		}

		result := get(l)
		out.WriteString(intToString(int64(result)))
		out.WriteByte('\n')
	}
}

func intToString(x int64) string {
	if x == 0 {
		return "0"
	}
	var buf [20]byte
	idx := len(buf)
	neg := x < 0
	if neg {
		x = -x
	}
	for x > 0 {
		idx--
		buf[idx] = byte('0' + x%10)
		x /= 10
	}
	if neg {
		idx--
		buf[idx] = '-'
	}
	return string(buf[idx:])
}
