package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const B = 20
const modVal = int64(1e9 + 7)

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) Add(i int, v int64) {
	for i <= b.n {
		b.tree[i] += v
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func (b *BIT) RangeSum(l, r int) int64 {
	if l > r {
		return 0
	}
	return b.Sum(r) - b.Sum(l-1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	rand.Seed(time.Now().UnixNano())

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	bits := make([]*BIT, B)
	for i := 0; i < B; i++ {
		bits[i] = NewBIT(n)
	}

	weights := make(map[int][]int64)
	getWeights := func(x int) []int64 {
		if w, ok := weights[x]; ok {
			return w
		}
		w := make([]int64, B)
		for i := 0; i < B; i++ {
			w[i] = rand.Int63n(modVal-1) + 1
		}
		weights[x] = w
		return w
	}

	for i := 1; i <= n; i++ {
		w := getWeights(arr[i])
		for j := 0; j < B; j++ {
			bits[j].Add(i, w[j])
		}
	}

	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(reader, &typ)
		if typ == 1 {
			var idx, x int
			fmt.Fscan(reader, &idx, &x)
			if arr[idx] == x {
				continue
			}
			wOld := getWeights(arr[idx])
			wNew := getWeights(x)
			for j := 0; j < B; j++ {
				bits[j].Add(idx, wNew[j]-wOld[j])
			}
			arr[idx] = x
		} else if typ == 2 {
			var l, r, k int
			fmt.Fscan(reader, &l, &r, &k)
			if k == 1 {
				fmt.Fprintln(writer, "YES")
				continue
			}
			ok := true
			kk := int64(k)
			for j := 0; j < B; j++ {
				sum := bits[j].RangeSum(l, r)
				if sum%kk != 0 {
					ok = false
					break
				}
			}
			if ok {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
