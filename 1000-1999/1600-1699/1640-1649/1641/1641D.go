package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

type subset struct {
	key  [5]int
	bits int
}

type item struct {
	w    int
	nums []int
	subs []subset
}

func makeSubsets(nums []int) []subset {
	m := len(nums)
	var res []subset
	for mask := 1; mask < 1<<m; mask++ {
		var k [5]int
		idx := 0
		for i := 0; i < m; i++ {
			if mask&(1<<i) != 0 {
				k[idx] = nums[i]
				idx++
			}
		}
		res = append(res, subset{key: k, bits: bits.OnesCount(uint(mask))})
	}
	return res
}

func addSubsets(cnt map[[5]int]int, it *item) {
	for _, s := range it.subs {
		cnt[s.key]++
	}
}

func removeSubsets(cnt map[[5]int]int, it *item) {
	for _, s := range it.subs {
		cnt[s.key]--
	}
}

func query(it *item, cnt map[[5]int]int, total int) bool {
	sum := 0
	for _, s := range it.subs {
		c := cnt[s.key]
		if s.bits%2 == 1 {
			sum += c
		} else {
			sum -= c
		}
	}
	return sum < total
}

func check(items []*item, weights []int, X int) bool {
	cnt := make(map[[5]int]int)
	p := 0
	n := len(items)
	for j := 0; j < n; j++ {
		wj := items[j].w
		if wj > X {
			break
		}
		limit := X - wj
		posLimit := sort.Search(n, func(i int) bool { return weights[i] > limit })
		if posLimit > j {
			posLimit = j
		}
		for p > posLimit {
			p--
			removeSubsets(cnt, items[p])
		}
		for p < posLimit {
			addSubsets(cnt, items[p])
			p++
		}
		if query(items[j], cnt, p) {
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	items := make([]*item, n)
	for i := 0; i < n; i++ {
		nums := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &nums[j])
		}
		var w int
		fmt.Fscan(in, &w)
		sort.Ints(nums)
		items[i] = &item{w: w, nums: nums}
		items[i].subs = makeSubsets(nums)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].w < items[j].w })
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		weights[i] = items[i].w
	}
	lo, hi := 2, items[n-1].w+items[n-2].w
	if hi < 2 {
		hi = 2
	}
	ans := -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if check(items, weights, mid) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	fmt.Println(ans)
}
