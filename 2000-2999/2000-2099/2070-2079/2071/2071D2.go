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
		r := fs.nextInt64()

		preXor := make([]int, n+1)
		sumA := make([]int64, n+1)
		prefPref := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			v := fs.nextInt()
			preXor[i] = preXor[i-1] ^ v
			sumA[i] = sumA[i-1] + int64(v)
			prefPref[i] = prefPref[i-1] + int64(preXor[i])
		}

		baseXor := preXor[n]
		if n%2 == 0 {
			baseXor ^= preXor[n/2]
		}
		nHalf := n / 2

		prefMemo := make(map[int64]int)
		sumPrefMemo := make(map[int64]int64)

		var prefX func(int64) int
		prefX = func(k int64) int {
			if k <= int64(n) {
				return preXor[int(k)]
			}
			if val, ok := prefMemo[k]; ok {
				return val
			}
			res := baseXor
			if k%2 == 0 {
				res ^= prefX(k / 2)
			}
			prefMemo[k] = res
			return res
		}

		var sumPref func(int64) int64
		sumPref = func(k int64) int64 {
			if k <= 0 {
				return 0
			}
			if k <= int64(n) {
				return prefPref[int(k)]
			}
			if val, ok := sumPrefMemo[k]; ok {
				return val
			}
			half := k / 2
			var res int64
			if baseXor == 0 {
				res = prefPref[n] + sumPref(half) - prefPref[nHalf]
			} else {
				res = prefPref[n] + (k - int64(n)) - sumPref(half) + prefPref[nHalf]
			}
			sumPrefMemo[k] = res
			return res
		}

		var sum func(int64) int64
		sum = func(k int64) int64 {
			if k <= 0 {
				return 0
			}
			if k <= int64(n) {
				return sumA[int(k)]
			}
			start := int64(0)
			if n%2 == 0 {
				start = int64(prefX(int64(nHalf)))
			}
			left := int64(nHalf + 1)
			right := (k - 1) / 2
			mid := int64(0)
			if right >= left {
				mid = sumPref(right) - sumPref(left-1)
			}
			end := int64(0)
			if k%2 == 0 {
				end = int64(prefX(k / 2))
			}
			return sumA[n] + start + 2*mid + end
		}

		ans := sum(r) - sum(l-1)
		out.WriteString(intToString(ans))
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
