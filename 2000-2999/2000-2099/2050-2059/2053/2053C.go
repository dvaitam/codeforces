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

func (fs *fastScanner) nextInt64() int64 {
	var sign int64 = 1
	var val int64
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b > '~') {
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

type result struct {
	sum int64
	cnt int64
}

func solve(len, k int64, memo map[int64]result) result {
	if len < k {
		return result{}
	}
	if val, ok := memo[len]; ok {
		return val
	}
	if len%2 == 0 {
		child := solve(len/2, k, memo)
		val := result{
			sum: child.sum*2 + (len/2)*child.cnt,
			cnt: child.cnt * 2,
		}
		memo[len] = val
		return val
	}
	child := solve(len/2, k, memo)
	mid := len/2 + 1
	val := result{
		sum: mid + child.sum*2 + mid*child.cnt,
		cnt: child.cnt*2 + 1,
	}
	memo[len] = val
	return val
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt64()
	for ; t > 0; t-- {
		n := fs.nextInt64()
		k := fs.nextInt64()
		memo := make(map[int64]result)
		ans := solve(n, k, memo).sum
		out.WriteString(int64ToString(ans))
		out.WriteByte('\n')
	}
}

func int64ToString(x int64) string {
	if x == 0 {
		return "0"
	}
	var buf [32]byte
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
