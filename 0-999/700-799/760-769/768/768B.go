package main

import (
	"bufio"
	"fmt"
	"os"
)

var lenMemo map[int64]int64
var onesMemo map[int64]int64

func seqLen(n int64) int64 {
	if n <= 1 {
		return 1
	}
	if v, ok := lenMemo[n]; ok {
		return v
	}
	v := 2*seqLen(n/2) + 1
	lenMemo[n] = v
	return v
}

func onesCount(n int64) int64 {
	if n <= 1 {
		return n
	}
	if v, ok := onesMemo[n]; ok {
		return v
	}
	v := 2*onesCount(n/2) + n%2
	onesMemo[n] = v
	return v
}

func prefix(n, pos int64) int64 {
	if pos <= 0 || n == 0 {
		return 0
	}
	if n == 1 {
		if pos >= 1 {
			return 1
		}
		return 0
	}
	leftLen := seqLen(n / 2)
	if pos <= leftLen {
		return prefix(n/2, pos)
	}
	leftOnes := onesCount(n / 2)
	if pos == leftLen+1 {
		return leftOnes + n%2
	}
	return leftOnes + n%2 + prefix(n/2, pos-leftLen-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, l, r int64
	if _, err := fmt.Fscan(in, &n, &l, &r); err != nil {
		return
	}
	lenMemo = make(map[int64]int64)
	onesMemo = make(map[int64]int64)
	ans := prefix(n, r) - prefix(n, l-1)
	fmt.Println(ans)
}
