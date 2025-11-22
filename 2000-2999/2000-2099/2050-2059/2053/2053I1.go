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

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min64(a, b int64) int64 {
	if a < b {
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
		n := int(in.nextInt64())
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			a[i] = in.nextInt64()
			sum += a[i]
		}

		var L, R int64 = 0, 0 // possible current prefix values without extra insertion
		ins := 0

		for i := 0; i < n; i++ {
			ai := a[i]
			JL := max64(0, -ai)
			JR := min64(sum, sum-ai)

			nL := max64(L, JL)
			nR := min64(R, JR)
			if nL > nR {
				ins++
				L, R = JL, JR
			} else {
				L, R = nL, nR
			}

			L += ai
			R += ai
		}

		if L > sum || R < sum {
			ins++
		}

		fmt.Fprintln(out, int64(n)+int64(ins))
	}
}
