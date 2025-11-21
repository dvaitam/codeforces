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

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		m := fs.nextInt()
		limit := n * m
		mask := make([]uint8, limit+1)
		used := make([]int, 0)

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				val := fs.nextInt()
				if mask[val] == 0 {
					used = append(used, val)
				}
				if (i+j)%2 == 0 {
					mask[val] |= 1
				} else {
					mask[val] |= 2
				}
			}
		}

		total := 0
		best := 0
		for _, val := range used {
			chi := 1
			if mask[val] == 3 {
				chi = 2
			}
			total += chi
			if chi > best {
				best = chi
			}
		}
		ans := total - best
		out.WriteString(intToString(ans))
		out.WriteByte('\n')
	}
}

func intToString(x int) string {
	if x == 0 {
		return "0"
	}
	var buf [16]byte
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
