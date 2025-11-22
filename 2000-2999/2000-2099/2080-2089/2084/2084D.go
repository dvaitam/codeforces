package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastReader struct {
	r *bufio.Reader
}

func newFastReader() *fastReader {
	return &fastReader{r: bufio.NewReader(os.Stdin)}
}

func (fr *fastReader) nextInt64() int64 {
	var sign int64 = 1
	var val int64
	c, err := fr.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fr.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, _ = fr.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fr.r.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := newFastReader()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(in.nextInt64())
	for ; t > 0; t-- {
		n := int64(in.nextInt64())
		m := int64(in.nextInt64())
		k := int64(in.nextInt64())

		// Maximum achievable mex
		x := min64(n/(m+1), n-m*k)
		d := max64(x, k)

		ans := make([]int64, n)
		fillVal := int64(1_000_000_000)
		for i := range ans {
			ans[i] = fillVal
		}

		for tBlock := int64(0); tBlock <= m; tBlock++ {
			start := tBlock * d
			for v := int64(0); v < x; v++ {
				ans[start+v] = v
			}
		}

		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
