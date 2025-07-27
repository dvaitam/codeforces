package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	val int
	idx int
}

type fenwickVal struct {
	n int
	t []pair
}

func newFenwickVal(n int) *fenwickVal {
	t := make([]pair, n+2)
	return &fenwickVal{n, t}
}
func (f *fenwickVal) update(i int, v pair) {
	for ; i <= f.n; i += i & -i {
		if v.val > f.t[i].val {
			f.t[i] = v
		}
	}
}
func (f *fenwickVal) query(i int) pair {
	res := pair{0, 0}
	for ; i > 0; i -= i & -i {
		if f.t[i].val > res.val {
			res = f.t[i]
		}
	}
	return res
}

type fenwickMin struct {
	n int
	t []int
}

func newFenwickMin(n int) *fenwickMin {
	t := make([]int, n+3)
	return &fenwickMin{n, t}
}
func (f *fenwickMin) update(i, v int) {
	for ; i <= f.n+1; i += i & -i {
		if v > f.t[i] {
			f.t[i] = v
		}
	}
}
func (f *fenwickMin) query(i int) int {
	res := 0
	for ; i > 0; i -= i & -i {
		if f.t[i] > res {
			res = f.t[i]
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		maxv := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] > maxv {
				maxv = a[i]
			}
		}
		fv := newFenwickVal(maxv)
		fm := newFenwickMin(maxv)
		ans := 0
		for _, x := range a {
			q1 := fv.query(x)
			q2 := fm.query(x + 1)
			pair1 := 0
			if q1.idx > 0 {
				pair1 = q1.idx
			}
			pair2 := x
			dp := 0
			pairMin := 0
			if q1.val > q2 || (q1.val == q2 && pair1 <= pair2) {
				dp = q1.val + 1
				pairMin = pair1
			} else {
				dp = q2 + 1
				pairMin = pair2
			}
			fv.update(x, pair{dp, x})
			fm.update(pairMin+1, dp)
			if dp > ans {
				ans = dp
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
