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

func (fs *fastScanner) nextString() string {
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b == '\n') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return ""
	}
	buf := []byte{b}
	for {
		b, err = fs.r.ReadByte()
		if err != nil || b <= ' ' || b == '\n' {
			break
		}
		buf = append(buf, b)
	}
	return string(buf)
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		_ = fs.nextInt() // n, not needed beyond bounds check
		s := fs.nextString()

		var cntMinus, cntUnder int64
		for i := 0; i < len(s); i++ {
			if s[i] == '-' {
				cntMinus++
			} else if s[i] == '_' {
				cntUnder++
			}
		}

		left := cntMinus / 2
		right := (cntMinus + 1) / 2
		ans := cntUnder * left * right
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
