package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) readByte() byte {
	b, err := fs.r.ReadByte()
	if err != nil {
		return 0
	}
	return b
}

func (fs *fastScanner) nextInt() int {
	c := fs.readByte()
	for c <= ' ' && c != 0 {
		c = fs.readByte()
	}
	sign := 1
	if c == '-' {
		sign = -1
		c = fs.readByte()
	}
	val := 0
	for c > ' ' {
		val = val*10 + int(c-'0')
		c = fs.readByte()
	}
	return sign * val
}

func (fs *fastScanner) nextString() string {
	c := fs.readByte()
	for c <= ' ' && c != 0 {
		c = fs.readByte()
	}
	if c == 0 {
		return ""
	}
	buf := make([]byte, 0, 16)
	for c > ' ' {
		buf = append(buf, c)
		c = fs.readByte()
	}
	return string(buf)
}

func reverseString(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		words := make([]string, n)
		counts := make([]int, 4) // 00, 11, 01, 10
		for i := 0; i < n; i++ {
			s := fs.nextString()
			words[i] = s
			first := s[0]
			last := s[len(s)-1]
			switch {
			case first == '0' && last == '0':
				counts[0]++
			case first == '1' && last == '1':
				counts[1]++
			case first == '0' && last == '1':
				counts[2]++
			default:
				counts[3]++
			}
		}

		mp := make(map[string]struct{}, n*2)
		for _, w := range words {
			mp[w] = struct{}{}
		}

		cand01 := make([]int, 0)
		cand10 := make([]int, 0)
		for idx, w := range words {
			first := w[0]
			last := w[len(w)-1]
			if first == last {
				continue
			}
			rev := reverseString(w)
			if _, exists := mp[rev]; exists {
				continue
			}
			if first == '0' {
				cand01 = append(cand01, idx+1)
			} else {
				cand10 = append(cand10, idx+1)
			}
		}

		c00, c11, c01, c10 := counts[0], counts[1], counts[2], counts[3]

		if c01+c10 == 0 {
			if c00 > 0 && c11 > 0 {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, 0)
			}
			continue
		}

		diff := c01 - c10
		switch {
		case diff > 1:
			need := diff / 2
			if len(cand01) < need {
				fmt.Fprintln(out, -1)
				continue
			}
			fmt.Fprintln(out, need)
			for i := 0; i < need; i++ {
				if i > 0 {
					out.WriteByte(' ')
				}
				fmt.Fprint(out, cand01[i])
			}
			out.WriteByte('\n')
		case diff < -1:
			need := (-diff) / 2
			if len(cand10) < need {
				fmt.Fprintln(out, -1)
				continue
			}
			fmt.Fprintln(out, need)
			for i := 0; i < need; i++ {
				if i > 0 {
					out.WriteByte(' ')
				}
				fmt.Fprint(out, cand10[i])
			}
			out.WriteByte('\n')
		default:
			fmt.Fprintln(out, 0)
		}
	}
}
