package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int = 998244353

type Seg struct {
	l int
	r int
}

func solveQuery(L, R int, segs []Seg) int {
	if len(segs) == 0 {
		return 0
	}
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].l == segs[j].l {
			return segs[i].r < segs[j].r
		}
		return segs[i].l < segs[j].l
	})
	n := len(segs)
	starts := make([]int, n)
	ends := make([]int, n)
	for i, s := range segs {
		starts[i] = s.l
		ends[i] = s.r
	}
	hi := make([]int, n)
	j := 0
	for i := 0; i < n; i++ {
		if j < i+1 {
			j = i + 1
		}
		for j < n && starts[j] <= ends[i] {
			j++
		}
		hi[i] = j - 1
	}
	diff := make([]int, n+1)
	pre := 0
	val := make([]int, n)
	res := 0
	for i := 0; i < n; i++ {
		pre += diff[i]
		pre %= mod
		if pre < 0 {
			pre += mod
		}
		val[i] = (val[i] + pre) % mod
		if starts[i] == L {
			val[i]--
			if val[i] < 0 {
				val[i] += mod
			}
		}
		if hi[i] >= i+1 {
			diff[i+1] = (diff[i+1] - val[i]) % mod
			diff[hi[i]+1] = (diff[hi[i]+1] + val[i]) % mod
		}
		if ends[i] == R {
			res += val[i]
			res %= mod
		}
	}
	if res < 0 {
		res += mod
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m, q int
	if _, err := fmt.Fscan(in, &m, &q); err != nil {
		return
	}
	segsAll := make([]Seg, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &segsAll[i].l, &segsAll[i].r)
	}
	sort.Slice(segsAll, func(i, j int) bool { return segsAll[i].l < segsAll[j].l })

	for ; q > 0; q-- {
		var L, R int
		fmt.Fscan(in, &L, &R)
		// binary search range of segments with start in [L, R]
		left := sort.Search(len(segsAll), func(i int) bool { return segsAll[i].l >= L })
		right := sort.Search(len(segsAll), func(i int) bool { return segsAll[i].l > R })
		tmp := make([]Seg, 0, right-left)
		for i := left; i < right; i++ {
			if segsAll[i].r <= R {
				tmp = append(tmp, segsAll[i])
			}
		}
		ans := solveQuery(L, R, tmp)
		fmt.Fprintln(out, ans)
	}
}
