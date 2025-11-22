package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) next() (uint64, bool) {
	var v uint64
	c, err := fs.r.ReadByte()
	for err == nil && (c == ' ' || c == '\n' || c == '\r' || c == '\t') {
		c, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0, false
	}
	for err == nil && c >= '0' && c <= '9' {
		v = v*10 + uint64(c-'0')
		c, err = fs.r.ReadByte()
	}
	if err == nil {
		fs.r.UnreadByte()
	}
	return v, true
}

func main() {
	in := newScanner()
	tVal, ok := in.next()
	if !ok {
		return
	}
	t := int(tVal)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; t > 0; t-- {
		nVal, _ := in.next()
		n := int(nVal)
		a := make([]uint64, n)
		for i := 0; i < n; i++ {
			x, _ := in.next()
			a[i] = x
		}

		left := make([]uint64, n)
		right := make([]uint64, n)

		left[0] = a[0]
		for i := 1; i < n; i++ {
			left[i] = a[i] + (left[i-1] >> 1)
		}

		right[n-1] = a[n-1]
		for i := n - 2; i >= 0; i-- {
			right[i] = a[i] + (right[i+1] >> 1)
		}

		for i := 0; i < n; i++ {
			res := left[i] + right[i] - a[i]
			if i+1 == n {
				fmt.Fprintln(out, res)
			} else {
				fmt.Fprint(out, res, " ")
			}
		}
	}
}
