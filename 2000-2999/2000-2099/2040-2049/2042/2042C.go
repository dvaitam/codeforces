package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		val = val*10 + int(b-'0')
		b, err = fs.r.ReadByte()
	}
	return sign * val
}

func (fs *fastScanner) nextString() string {
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b > '~') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return ""
	}
	buf := []byte{b}
	for {
		b, err = fs.r.ReadByte()
		if err != nil || b <= ' ' {
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
		n := fs.nextInt()
		k := int64(fs.nextInt())
		s := fs.nextString()

		suff := make([]int64, n+1)
		for i := n - 1; i >= 0; i-- {
			val := int64(1)
			if s[i] == '0' {
				val = -1
			}
			suff[i] = suff[i+1] + val
		}

		vals := make([]int64, 0, n-1)
		for i := 1; i < n; i++ {
			vals = append(vals, suff[i])
		}
		sort.Slice(vals, func(i, j int) bool {
			return vals[i] > vals[j]
		})

		pref := make([]int64, len(vals)+1)
		for i, v := range vals {
			pref[i+1] = pref[i] + v
		}

		answer := -1
		for m := 1; m <= n; m++ {
			if pref[m-1] >= k {
				answer = m
				break
			}
		}
		fmt.Fprintln(out, answer)
	}
}
