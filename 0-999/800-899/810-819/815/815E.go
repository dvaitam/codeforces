package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	length int64
	d      int64
}

var memoGeq map[pair]int64

func countGeq(length, d int64) int64 {
	if length <= 0 {
		return 0
	}
	key := pair{length, d}
	if v, ok := memoGeq[key]; ok {
		return v
	}
	dist := (length + 1) / 2
	if dist < d {
		memoGeq[key] = 0
		return 0
	}
	left := (length - 1) / 2
	right := length - 1 - left
	res := int64(1) + countGeq(left, d) + countGeq(right, d)
	memoGeq[key] = res
	return res
}

func countEq(length, d int64) int64 {
	return countGeq(length, d) - countGeq(length, d+1)
}

func kthWithDist(length, start, d, t int64) int64 {
	dist := (length + 1) / 2
	leftLen := (length - 1) / 2
	rightLen := length - 1 - leftLen
	mid := start + leftLen
	if dist < d {
		return -1
	}
	if dist == d {
		leftEq := countEq(leftLen, d)
		if t <= leftEq {
			return kthWithDist(leftLen, start, d, t)
		}
		if t == leftEq+1 {
			return mid
		}
		return kthWithDist(rightLen, mid+1, d, t-leftEq-1)
	}
	leftEq := countEq(leftLen, d)
	if t <= leftEq {
		return kthWithDist(leftLen, start, d, t)
	}
	return kthWithDist(rightLen, mid+1, d, t-leftEq)
}

func solve(n, k int64) int64 {
	if k == 1 {
		return 1
	}
	if k == 2 {
		return n
	}
	length := n - 2
	lo, hi := int64(1), (length+1)/2
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if countGeq(length, mid) >= k-2 {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	d := lo
	t := k - 2 - countGeq(length, d+1)
	return kthWithDist(length, 2, d, t)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	memoGeq = make(map[pair]int64)

	ans := solve(n, k)
	fmt.Fprintln(writer, ans)
}
