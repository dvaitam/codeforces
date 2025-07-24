package main

import (
	"bufio"
	"fmt"
	"os"
)

func count(n, y int64) int64 {
	var length int64 = 1
	if y%2 == 0 {
		length = 2
	}
	var res int64
	cur := y
	l := length
	for cur <= n {
		end := cur + l - 1
		if end > n {
			res += n - cur + 1
		} else {
			res += l
		}
		if cur > n/2 {
			break
		}
		cur <<= 1
		l <<= 1
	}
	return res
}

func maxY(n, k int64, parity int64) int64 {
	var low int64
	if parity == 1 {
		low = 1
	} else {
		low = 2
	}
	high := n
	if high%2 != parity {
		high--
	}
	var ans int64
	for low <= high {
		mid := (low + high) / 2
		if mid%2 != parity {
			mid++
			if mid > high {
				break
			}
		}
		if count(n, mid) >= k {
			if mid > ans {
				ans = mid
			}
			low = mid + 2
		} else {
			high = mid - 2
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	oddAns := maxY(n, k, 1)
	evenAns := maxY(n, k, 0)
	if oddAns > evenAns {
		fmt.Println(oddAns)
	} else {
		fmt.Println(evenAns)
	}
}
