package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const prec = 1e9

type plane struct {
	a int64
	b int64
}

type BIT struct {
	n    int
	tree []int64
}

func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) add(i int, val int64) {
	for i <= b.n {
		b.tree[i] += val
		i += i & -i
	}
}

func (b *BIT) sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, w int
	if _, err := fmt.Fscan(in, &n, &w); err != nil {
		return
	}

	planes := make([]plane, n)
	for i := 0; i < n; i++ {
		var x, v int
		fmt.Fscan(in, &x, &v)
		a := -float64(x) / float64(v-w)
		b := -float64(x) / float64(v+w)
		planes[i] = plane{int64(math.Round(a * prec)), int64(math.Round(b * prec))}
	}

	sort.Slice(planes, func(i, j int) bool { return planes[i].a < planes[j].a })

	bvals := make([]int64, n)
	for i := range planes {
		bvals[i] = planes[i].b
	}
	uniqMap := map[int64]struct{}{}
	uniq := make([]int64, 0, len(bvals))
	for _, v := range bvals {
		if _, ok := uniqMap[v]; !ok {
			uniqMap[v] = struct{}{}
			uniq = append(uniq, v)
		}
	}
	sort.Slice(uniq, func(i, j int) bool { return uniq[i] < uniq[j] })
	rank := make(map[int64]int, len(uniq))
	for i, v := range uniq {
		rank[v] = i + 1
	}

	bit := newBIT(len(uniq) + 2)
	var total int64
	var ans int64
	i := 0
	for i < n {
		j := i
		for j < n && planes[j].a == planes[i].a {
			j++
		}
		for k := i; k < j; k++ {
			r := rank[planes[k].b]
			ans += total - bit.sum(r-1)
		}
		for k := i; k < j; k++ {
			r := rank[planes[k].b]
			bit.add(r, 1)
			total++
		}
		m := int64(j - i)
		ans += m * (m - 1) / 2
		i = j
	}

	fmt.Println(ans)
}
