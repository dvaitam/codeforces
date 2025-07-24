package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Item struct {
	cost  int
	group int
	idx   int
}

type Pair struct {
	cost int
	pos  int
}

type BIT struct {
	n   int
	cnt []int
	sum []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, cnt: make([]int, n+2), sum: make([]int64, n+2)}
}

func (b *BIT) Add(pos, dc int, ds int64) {
	for i := pos; i <= b.n; i += i & -i {
		b.cnt[i] += dc
		b.sum[i] += ds
	}
}

func (b *BIT) PrefixCnt(pos int) int {
	s := 0
	for i := pos; i > 0; i &= i - 1 {
		s += b.cnt[i]
	}
	return s
}

func (b *BIT) PrefixSum(pos int) int64 {
	var s int64
	for i := pos; i > 0; i &= i - 1 {
		s += b.sum[i]
	}
	return s
}

func (b *BIT) Total() int { return b.PrefixCnt(b.n) }

func (b *BIT) Kth(k int) int {
	pos := 0
	bit := 1
	for bit<<1 <= b.n {
		bit <<= 1
	}
	for bit > 0 {
		nxt := pos + bit
		if nxt <= b.n && b.cnt[nxt] < k {
			k -= b.cnt[nxt]
			pos = nxt
		}
		bit >>= 1
	}
	return pos + 1
}

func (b *BIT) SumK(k int) int64 {
	if k <= 0 {
		return 0
	}
	if k > b.Total() {
		const INF = int64(1 << 60)
		return INF
	}
	idx := b.Kth(k)
	return b.PrefixSum(idx)
}

func prefix(arr []Pair) []int64 {
	pre := make([]int64, len(arr)+1)
	for i, p := range arr {
		pre[i+1] = pre[i] + int64(p.cost)
	}
	return pre
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	costs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &costs[i])
	}
	var a int
	fmt.Fscan(in, &a)
	likeM := make([]bool, n)
	for i := 0; i < a; i++ {
		var x int
		fmt.Fscan(in, &x)
		likeM[x-1] = true
	}
	var b int
	fmt.Fscan(in, &b)
	likeA := make([]bool, n)
	for i := 0; i < b; i++ {
		var x int
		fmt.Fscan(in, &x)
		likeA[x-1] = true
	}

	items := make([]Item, n)
	for i := 0; i < n; i++ {
		g := 4
		if likeM[i] && likeA[i] {
			g = 1
		} else if likeM[i] {
			g = 2
		} else if likeA[i] {
			g = 3
		}
		items[i] = Item{costs[i], g, i}
	}

	sortedItems := append([]Item(nil), items...)
	sort.Slice(sortedItems, func(i, j int) bool {
		if sortedItems[i].cost == sortedItems[j].cost {
			return sortedItems[i].idx < sortedItems[j].idx
		}
		return sortedItems[i].cost < sortedItems[j].cost
	})

	posMap := make([]int, n)
	costAtPos := make([]int, n+1)
	for i, it := range sortedItems {
		pos := i + 1
		posMap[it.idx] = pos
		costAtPos[pos] = it.cost
	}

	var both, mOnly, aOnly, none []Pair
	for _, it := range sortedItems {
		pair := Pair{it.cost, posMap[it.idx]}
		switch it.group {
		case 1:
			both = append(both, pair)
		case 2:
			mOnly = append(mOnly, pair)
		case 3:
			aOnly = append(aOnly, pair)
		default:
			none = append(none, pair)
		}
	}

	prefBoth := prefix(both)
	prefM := prefix(mOnly)
	prefA := prefix(aOnly)

	bit := NewBIT(n)
	for i := 1; i <= n; i++ {
		bit.Add(i, 1, int64(costAtPos[i]))
	}

	remB, remM, remA := 0, 0, 0
	maxX := min(len(both), m)
	ans := int64(1 << 60)

	for x := 0; x <= maxX; x++ {
		need := k - x
		if need < 0 {
			need = 0
		}
		targetM := min(len(mOnly), need)
		targetA := min(len(aOnly), need)

		for remB < x {
			bit.Add(both[remB].pos, -1, -int64(both[remB].cost))
			remB++
		}
		for remB > x {
			remB--
			bit.Add(both[remB].pos, 1, int64(both[remB].cost))
		}
		for remM < targetM {
			bit.Add(mOnly[remM].pos, -1, -int64(mOnly[remM].cost))
			remM++
		}
		for remM > targetM {
			remM--
			bit.Add(mOnly[remM].pos, 1, int64(mOnly[remM].cost))
		}
		for remA < targetA {
			bit.Add(aOnly[remA].pos, -1, -int64(aOnly[remA].cost))
			remA++
		}
		for remA > targetA {
			remA--
			bit.Add(aOnly[remA].pos, 1, int64(aOnly[remA].cost))
		}

		if need > len(mOnly) || need > len(aOnly) || x > len(both) {
			continue
		}
		r := m - (x + 2*need)
		if r < 0 || r > bit.Total() {
			continue
		}
		cost := prefBoth[x] + prefM[need] + prefA[need]
		other := bit.SumK(r)
		if other >= int64(1<<60) {
			continue
		}
		cost += other
		if cost < ans {
			ans = cost
		}
	}

	if ans == int64(1<<60) {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
