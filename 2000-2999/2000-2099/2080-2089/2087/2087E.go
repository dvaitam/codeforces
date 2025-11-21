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

	t := int(fs.nextInt64())
	for ; t > 0; t-- {
		n := int(fs.nextInt64())
		s := fs.nextString()
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			c[i] = fs.nextInt64()
		}

		pos := make([]int64, n)
		for i := 0; i < n; i++ {
			if c[i] > 0 {
				pos[i] = c[i]
			}
		}

		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + pos[i]
		}
		suffix := make([]int64, n+1)
		for i := n - 1; i >= 0; i-- {
			suffix[i] = suffix[i+1] + pos[i]
		}

		var best int64
		for i := 0; i < n; i++ {
			if s[i] == '>' {
				cur := c[i] + suffix[i+1]
				if cur > best {
					best = cur
				}
			} else {
				cur := c[i] + prefix[i]
				if cur > best {
					best = cur
				}
			}
		}
		if best < 0 {
			best = 0
		}
		out.WriteString(int64ToString(best))
		out.WriteByte('\n')
	}
}

func int64ToString(x int64) string {
	if x == 0 {
		return "0"
	}
	var buf [20]byte
	idx := len(buf)
	sign := x < 0
	if sign {
		x = -x
	}
	for x > 0 {
		idx--
		buf[idx] = byte('0' + x%10)
		x /= 10
	}
	if sign {
		idx--
		buf[idx] = '-'
	}
	return string(buf[idx:])
}
