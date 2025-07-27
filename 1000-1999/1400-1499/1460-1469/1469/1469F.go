package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func can(R int64, lengths []int64, k int64) bool {
	if R == 0 {
		return k <= 1
	}
	n := len(lengths)
	arr := make([]int64, n)
	copy(arr, lengths)
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	cnt := make([]int64, R+1)
	cnt[0] = 1
	total := int64(1)
	idx := int64(0)
	for _, L := range arr {
		for idx < R {
			if cnt[idx] == 0 {
				idx++
				continue
			}
			m := R - idx - 1
			if m <= 0 {
				idx++
				continue
			}
			break
		}
		if idx >= R || total >= k {
			break
		}
		d := idx
		cnt[d]--
		total--
		m := R - d - 1
		if L <= 2*m+1 {
			x := (L - 1) / 2
			for i := int64(1); i <= x && d+1+i <= R; i++ {
				cnt[d+1+i] += 2
			}
			if L%2 == 0 && d+1+L/2 <= R {
				cnt[d+1+L/2]++
			}
			total += L - 1
		} else {
			for i := int64(1); i <= m && d+1+i <= R; i++ {
				cnt[d+1+i] += 2
			}
			total += 2 * m
		}
		if total >= k {
			return true
		}
	}
	return total >= k
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k int64
	fmt.Fscan(in, &n, &k)
	lengths := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &lengths[i])
	}
	sum := int64(1)
	for _, L := range lengths {
		sum += L - 2
	}
	if sum < k {
		fmt.Fprintln(out, -1)
		return
	}
	l, r := int64(1), int64(200000)
	for l < r {
		mid := (l + r) / 2
		if can(mid, lengths, k) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	fmt.Fprintln(out, l)
}
