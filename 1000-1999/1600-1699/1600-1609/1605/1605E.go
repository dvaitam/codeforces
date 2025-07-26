package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type prefixHelper struct {
	arr  []int64
	pref []int64
}

func newPrefix(arr []int64) *prefixHelper {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	pref := make([]int64, len(arr)+1)
	for i, v := range arr {
		pref[i+1] = pref[i] + v
	}
	return &prefixHelper{arr, pref}
}

// sumPlus computes sum |x + t| over sorted arr
func (p *prefixHelper) sumPlus(t int64) int64 {
	if len(p.arr) == 0 {
		return 0
	}
	// find first index with arr[idx] >= -t
	idx := sort.Search(len(p.arr), func(i int) bool { return p.arr[i] >= -t })
	total := p.pref[len(p.arr)]
	left := p.pref[idx]
	return (total - 2*left) + int64(len(p.arr)-2*idx)*t
}

// sumMinus computes sum |t - x| over sorted arr
func (p *prefixHelper) sumMinus(t int64) int64 {
	if len(p.arr) == 0 {
		return 0
	}
	idx := sort.Search(len(p.arr), func(i int) bool { return p.arr[i] >= t })
	total := p.pref[len(p.arr)]
	left := p.pref[idx]
	return (total - 2*left) + int64(2*idx-len(p.arr))*t
}

func mobius(n int) []int {
	mu := make([]int, n+1)
	prime := make([]int, 0)
	isComp := make([]bool, n+1)
	mu[1] = 1
	for i := 2; i <= n; i++ {
		if !isComp[i] {
			prime = append(prime, i)
			mu[i] = -1
		}
		for _, p := range prime {
			if i*p > n {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			} else {
				mu[i*p] = -mu[i]
			}
		}
	}
	return mu
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}

	mu := mobius(n)

	diffBase := make([]int64, n+1)
	for i := 2; i <= n; i++ {
		diffBase[i] = b[i] - a[i]
	}

	base := make([]int64, n+1)
	for d := 1; d <= n; d++ {
		if mu[d] == 0 {
			continue
		}
		for m := d; m <= n; m += d {
			base[m] += int64(mu[d]) * diffBase[m/d]
		}
	}

	var plusArr, minusArr []int64
	constArrSum := int64(0)
	for i := 1; i <= n; i++ {
		if mu[i] == 1 {
			plusArr = append(plusArr, base[i])
		} else if mu[i] == -1 {
			minusArr = append(minusArr, base[i])
		} else {
			constArrSum += abs64(base[i])
		}
	}

	plusHelper := newPrefix(plusArr)
	minusHelper := newPrefix(minusArr)

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var x int64
		fmt.Fscan(in, &x)
		t := x - a[1]
		ans := constArrSum + plusHelper.sumPlus(t) + minusHelper.sumMinus(t)
		fmt.Fprintln(out, ans)
	}
}
