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

func (fs *fastScanner) nextInt64() int64 {
	var sign int64 = 1
	var val int64
	b, err := fs.r.ReadByte()
	for (b < '0' || b > '9') && b != '-' {
		if err != nil {
			return 0
		}
		b, err = fs.r.ReadByte()
	}
	if b == '-' {
		sign = -1
		b, err = fs.r.ReadByte()
	}
	for b >= '0' && b <= '9' {
		val = val*10 + int64(b-'0')
		b, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func singleValue(x int64) int64 {
	cnt := int64(1)
	for x > 3 {
		x = x/2 + 1
		cnt++
	}
	return cnt
}

func isSpecial(x int64) int64 {
	if x <= 2 {
		return 0
	}
	y := x - 1
	if y&(y-1) == 0 {
		return 1
	}
	return 0
}

func main() {
	fs := newScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(fs.nextInt64())
	for ; t > 0; t-- {
		n := int(fs.nextInt64())
		q := int(fs.nextInt64())
		prefSum := make([]int64, n+1)
		prefSpecial := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			val := fs.nextInt64()
			prefSum[i] = prefSum[i-1] + singleValue(val)
			prefSpecial[i] = prefSpecial[i-1] + isSpecial(val)
		}
		for ; q > 0; q-- {
			l := int(fs.nextInt64())
			r := int(fs.nextInt64())
			base := prefSum[r] - prefSum[l-1]
			cnt := prefSpecial[r] - prefSpecial[l-1]
			fmt.Fprintln(out, base+cnt/2)
		}
	}
}
