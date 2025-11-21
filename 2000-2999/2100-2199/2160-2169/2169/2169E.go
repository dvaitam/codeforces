package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	fullMask = 15
	negInf   = int64(-1 << 60)
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
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func bitPos(lsb int) int {
	switch lsb {
	case 1:
		return 0
	case 2:
		return 1
	case 4:
		return 2
	default:
		return 3
	}
}

func main() {
	sc := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(sc.nextInt64())
	for ; t > 0; t-- {
		n := int(sc.nextInt64())
		xs := make([]int64, n)
		ys := make([]int64, n)
		cs := make([]int64, n)
		for i := 0; i < n; i++ {
			xs[i] = sc.nextInt64()
		}
		for i := 0; i < n; i++ {
			ys[i] = sc.nextInt64()
		}
		for i := 0; i < n; i++ {
			cs[i] = sc.nextInt64()
		}

		totalC := int64(0)
		dp := [fullMask + 1]int64{}
		for i := 1; i <= fullMask; i++ {
			dp[i] = negInf
		}

		for i := 0; i < n; i++ {
			totalC += cs[i]
			features := [4]int64{xs[i], -xs[i], ys[i], -ys[i]}
			var sums [fullMask + 1]int64
			var vals [fullMask + 1]int64
			for mask := 1; mask <= fullMask; mask++ {
				lsb := mask & -mask
				idx := bitPos(lsb)
				sums[mask] = sums[mask^lsb] + features[idx]
				vals[mask] = 2*sums[mask] - cs[i]
			}
			for s := fullMask; s >= 0; s-- {
				cur := dp[s]
				if cur == negInf {
					continue
				}
				rem := (^s) & fullMask
				sub := rem
				for sub > 0 {
					nm := s | sub
					cand := cur + vals[sub]
					if cand > dp[nm] {
						dp[nm] = cand
					}
					sub = (sub - 1) & rem
				}
			}
		}

		fmt.Fprintln(out, totalC+dp[fullMask])
	}
}
