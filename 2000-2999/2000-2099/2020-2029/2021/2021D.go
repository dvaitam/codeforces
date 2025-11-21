package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf int64 = -1 << 60

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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := newScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(in.nextInt64())
	for ; t > 0; t-- {
		n := int(in.nextInt64())
		m := int(in.nextInt64())

		row := make([]int64, m)
		readRow := func() {
			for i := 0; i < m; i++ {
				row[i] = in.nextInt64()
			}
		}

		// read first day and build base dp arrays
		readRow()
		pref := make([]int64, m+1)
		for i := 1; i <= m; i++ {
			pref[i] = pref[i-1] + row[i-1]
		}

		preMin := make([]int64, m+1)
		for i := 1; i <= m; i++ {
			preMin[i] = min(preMin[i-1], pref[i])
		}

		sufMax := make([]int64, m+2)
		sufMax[m] = pref[m]
		for i := m - 1; i >= 0; i-- {
			sufMax[i] = max(sufMax[i+1], pref[i])
		}

		prevL := make([]int64, m+1) // best total with interval starting at i
		prevR := make([]int64, m+1) // best total with interval ending at i
		for l := 1; l <= m; l++ {
			prevL[l] = sufMax[l] - pref[l-1]
		}
		for r := 1; r <= m; r++ {
			prevR[r] = pref[r] - preMin[r-1]
		}

		// process remaining days
		for day := 1; day < n; day++ {
			readRow()

			for i := 1; i <= m; i++ {
				pref[i] = pref[i-1] + row[i-1]
			}
			preMin[0] = 0
			for i := 1; i <= m; i++ {
				preMin[i] = min(preMin[i-1], pref[i])
			}
			sufMax[m] = pref[m]
			for i := m - 1; i >= 0; i-- {
				sufMax[i] = max(sufMax[i+1], pref[i])
			}

			curR := make([]int64, m+1)
			leftRun, rightRun := negInf, negInf
			for x := 1; x <= m; x++ {
				if x >= 2 {
					leftRun = max(leftRun, prevL[x]-preMin[x-2]) // need a new element strictly to the left
				}
				bestPrev := max(leftRun, rightRun)
				curR[x] = pref[x] + bestPrev
				rightRun = max(rightRun, prevR[x]-preMin[x-1]) // can reuse from right endpoint as long as we start before it
			}

			curL := make([]int64, m+1)
			run1, run2 := negInf, negInf
			for l := m; l >= 1; l-- {
				if l+1 <= m {
					run1 = max(run1, prevL[l+1]+sufMax[l+1]) // extend to the left of previous left endpoint
				}
				if l <= m-1 {
					run2 = max(run2, prevR[l]+sufMax[l+1]) // extend to the right of previous right endpoint
				}
				curL[l] = -pref[l-1] + max(run1, run2)
			}

			prevL, prevR = curL, curR
		}

		ans := negInf
		for i := 1; i <= m; i++ {
			if prevL[i] > ans {
				ans = prevL[i]
			}
			if prevR[i] > ans {
				ans = prevR[i]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
