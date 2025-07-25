package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(a, b, r int64) int64 {
	const maxBit = 60
	var w [maxBit + 1]int64
	var absW [maxBit + 1]int64
	for i := 0; i <= maxBit; i++ {
		aBit := (a >> i) & 1
		bBit := (b >> i) & 1
		if aBit != bBit {
			if aBit > bBit {
				w[i] = 1 << (i + 1)
			} else {
				w[i] = -1 << (i + 1)
			}
		} else {
			w[i] = 0
		}
		if w[i] >= 0 {
			absW[i] = w[i]
		} else {
			absW[i] = -w[i]
		}
	}

	var prefAll [maxBit + 2]int64
	var prefLim [maxBit + 2]int64
	for i := 1; i <= maxBit+1; i++ {
		prefAll[i] = prefAll[i-1] + absW[i-1]
		if (r>>(i-1))&1 == 1 {
			prefLim[i] = prefLim[i-1] + absW[i-1]
		} else {
			prefLim[i] = prefLim[i-1]
		}
	}

	delta := a - b
	var x int64
	prefixLess := false
	for i := maxBit; i >= 0; i-- {
		remAll := prefAll[i]
		remLim := prefLim[i]
		if prefixLess {
			delta0 := delta
			best0 := abs(delta0) - remAll
			if best0 < 0 {
				best0 = 0
			}
			delta1 := delta - w[i]
			best1 := abs(delta1) - remAll
			if best1 < 0 {
				best1 = 0
			}
			if best1 < best0 || (best1 == best0 && abs(delta1) < abs(delta0)) {
				delta = delta1
				x |= 1 << i
			}
		} else {
			if (r>>i)&1 == 0 {
				continue
			}
			deltaEqual := delta - w[i]
			bestEqual := abs(deltaEqual) - remLim
			if bestEqual < 0 {
				bestEqual = 0
			}
			deltaLower := delta
			bestLower := abs(deltaLower) - remAll
			if bestLower < 0 {
				bestLower = 0
			}
			if bestLower < bestEqual || (bestLower == bestEqual && abs(deltaLower) < abs(deltaEqual)) {
				prefixLess = true
			} else {
				delta = deltaEqual
				x |= 1 << i
			}
		}
	}
	if delta < 0 {
		delta = -delta
	}
	return delta
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b, r int64
		fmt.Fscan(reader, &a, &b, &r)
		ans := solveCase(a, b, r)
		fmt.Fprintln(writer, ans)
	}
}
