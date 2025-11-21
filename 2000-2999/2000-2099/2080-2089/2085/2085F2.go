package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

type Fenwick struct {
	n      int
	bit    []int64
	maxPow int
}

func NewFenwick(size int) *Fenwick {
	bit := make([]int64, size+2)
	maxPow := 1
	for maxPow < size+2 {
		maxPow <<= 1
	}
	return &Fenwick{n: size + 1, bit: bit, maxPow: maxPow}
}

func (f *Fenwick) Add(idx int, delta int64) {
	if idx <= 0 {
		return
	}
	for idx <= f.n {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}

func (f *Fenwick) Sum(idx int) int64 {
	if idx > f.n {
		idx = f.n
	}
	res := int64(0)
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func (f *Fenwick) Kth(k int64) int {
	idx := 0
	bitMask := f.maxPow
	for bitMask > 0 {
		next := idx + bitMask
		if next <= f.n && f.bit[next] < k {
			k -= f.bit[next]
			idx = next
		}
		bitMask >>= 1
	}
	return idx + 1
}

type CostState struct {
	L             int
	pos           [][]int
	occIdx        []int
	bitCnt, bitSum *Fenwick
	needLeftCnt   int
	needRightCnt  int
	countBoth     int
	needLeftSum   int64
	needRightSum  int64
	sumRightBoth  int64
	sumValuesBoth int64
	maxVal        int
}

func NewCostState(k, maxPos int, pos [][]int) *CostState {
	maxVal := 2*maxPos + 5
	return &CostState{
		L:             k / 2,
		pos:           pos,
		occIdx:        make([]int, len(pos)),
		bitCnt:        NewFenwick(maxVal),
		bitSum:        NewFenwick(maxVal),
		maxVal:        maxVal,
	}
}

func (cs *CostState) leftVal(c int) int {
	idx := cs.occIdx[c]
	return cs.pos[c][idx-1]
}

func (cs *CostState) rightVal(c int) int {
	idx := cs.occIdx[c]
	return cs.pos[c][idx]
}

func (cs *CostState) removeColor(c int) {
	idx := cs.occIdx[c]
	list := cs.pos[c]
	if len(list) == 0 {
		return
	}
	hasLeft := idx > 0
	hasRight := idx < len(list)
	switch {
	case hasLeft && hasRight:
		lv := cs.leftVal(c)
		rv := cs.rightVal(c)
		cs.countBoth--
		cs.sumRightBoth -= int64(rv)
		cs.sumValuesBoth -= int64(lv + rv)
		cs.bitCnt.Add(lv+rv, -1)
		cs.bitSum.Add(lv+rv, -int64(lv+rv))
	case hasRight:
		rv := cs.rightVal(c)
		cs.needRightCnt--
		cs.needRightSum -= int64(rv)
	case hasLeft:
		lv := cs.leftVal(c)
		cs.needLeftCnt--
		cs.needLeftSum -= int64(lv)
	default:
	}
}

func (cs *CostState) addColor(c int) {
	idx := cs.occIdx[c]
	list := cs.pos[c]
	if len(list) == 0 {
		return
	}
	hasLeft := idx > 0
	hasRight := idx < len(list)
	switch {
	case hasLeft && hasRight:
		lv := cs.leftVal(c)
		rv := cs.rightVal(c)
		cs.countBoth++
		cs.sumRightBoth += int64(rv)
		cs.sumValuesBoth += int64(lv + rv)
		cs.bitCnt.Add(lv+rv, 1)
		cs.bitSum.Add(lv+rv, int64(lv+rv))
	case hasRight:
		rv := cs.rightVal(c)
		cs.needRightCnt++
		cs.needRightSum += int64(rv)
	case hasLeft:
		lv := cs.leftVal(c)
		cs.needLeftCnt++
		cs.needLeftSum += int64(lv)
	default:
	}
}

func (cs *CostState) sumSmallest(m int) int64 {
	if m <= 0 {
		return 0
	}
	pos := cs.bitCnt.Kth(int64(m))
	sum := cs.bitSum.Sum(pos)
	cnt := cs.bitCnt.Sum(pos)
	extra := cnt - int64(m)
	if extra > 0 {
		sum -= int64(pos) * extra
	}
	return sum
}

func (cs *CostState) sumLargest(k int) int64 {
	if k <= 0 {
		return 0
	}
	if k >= cs.countBoth {
		return cs.sumValuesBoth
	}
	small := cs.countBoth - k
	sumSmall := cs.sumSmallest(small)
	return cs.sumValuesBoth - sumSmall
}

func (cs *CostState) computeCost() (int64, bool) {
	if cs.needLeftCnt > cs.L || cs.needRightCnt > cs.L {
		return 0, false
	}
	chooseLeft := cs.L - cs.needLeftCnt
	if chooseLeft < 0 || chooseLeft > cs.countBoth {
		return 0, false
	}
	sumTop := cs.sumLargest(chooseLeft)
	cost := cs.needRightSum + cs.sumRightBoth - cs.needLeftSum - sumTop
	return cost, true
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		k := in.NextInt()
		arr := make([]int, n+1)
		pos := make([][]int, k+1)
		for i := 1; i <= n; i++ {
			v := in.NextInt()
			arr[i] = v
			pos[v] = append(pos[v], i)
		}

		cs := NewCostState(k, n, pos)
		for v := 1; v <= k; v++ {
			if len(pos[v]) == 0 {
				continue
			}
			cs.addColor(v)
		}

		const inf int64 = 1 << 60
		ans := inf

		L := k / 2
		var constant int64
		if k%2 == 0 {
			constant = int64(L * L)
		} else {
			constant = int64(L*L + L)
		}

		evaluate := func() {
			if cost, ok := cs.computeCost(); ok {
				val := cost - constant
				if val < ans {
					ans = val
				}
			}
		}

		kIsOdd := k%2 == 1

		for i := 1; i <= n; i++ {
			c := arr[i]
			cs.removeColor(c)
			cs.occIdx[c]++
			cs.addColor(c)
			evaluate()
			if kIsOdd {
				cs.removeColor(c)
				evaluate()
				cs.addColor(c)
			}
		}

		fmt.Fprintln(out, ans)
	}
}
