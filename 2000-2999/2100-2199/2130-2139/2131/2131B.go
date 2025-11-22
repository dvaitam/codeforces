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
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

func main() {
	in := newScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(in.nextInt64())
	for ; t > 0; t-- {
		n := int(in.nextInt64())
		ans := make([]int64, n)

		// pPrev: prefix sum up to i-1
		// bestPrev: maximum prefix sum up to i-1
		// bestPrevPrev: maximum prefix sum up to i-2
		var pPrev int64 = 0
		var bestPrev int64 = 0
		var bestPrevPrev int64 = 0

		for i := 1; i <= n; i++ {
			if i == 1 {
				ans[0] = -1
				pPrev = -1
				bestPrev = 0 // max of {0, -1}
				bestPrevPrev = 0
				continue
			}

			if i%2 == 0 { // even position, positive value
				need := bestPrevPrev + 1 // must exceed best prefix up to i-2
				if i < n {
					needNext := bestPrev + 2 // ensure next (odd) step can be -1
					if needNext > need {
						need = needNext
					}
				}
				val := need
				ans[i-1] = val - pPrev
				pPrev = val
				bestPrevPrev = bestPrev
				if pPrev > bestPrev {
					bestPrev = pPrev
				}
			} else { // odd position (>1), negative value
				need := bestPrevPrev + 1
				diff := pPrev - need // this is the absolute value placed with negative sign
				ans[i-1] = -diff
				pPrev = need
				bestPrevPrev = bestPrev
				// bestPrev unchanged as prefix sum decreased
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
