package main

import (
	"bufio"
	"fmt"
	"os"
)

const P = 998244353

func fpow(a, b int) int {
	res := 1
	for b > 1 {
		if b&1 == 1 {
			res = res * a % P
		}
		b >>= 1
		a = a * a % P
	}
	return a * res % P
}

// Fenwick tree for point updates and prefix sums
type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Add(i, v int) {
	for ; i <= f.n; i += i & -i {
		f.tree[i] += v
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += f.tree[i]
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n+1)
	vis := make([]bool, n+1)
	tot := 0
	fac := 1
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] == -1 {
			tot++
			fac = fac * tot % P
		} else {
			if a[i] >= 1 && a[i] <= n {
				vis[a[i]] = true
			}
		}
	}
	// initial ans: pairs of -1 expectation
	ans := tot * (tot - 1) % P * fpow(4, P-2) % P
	cnt := tot
	// count existing inversions
	fenw := NewFenwick(n)
	for i := n; i >= 1; i-- {
		if a[i] != -1 {
			ans = (ans + fenw.Sum(a[i]-1)) % P
			fenw.Add(a[i], 1)
		}
	}
	// prefix sum of missing values
	sum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		sum[i] = sum[i-1]
		if !vis[i] {
			sum[i]++
		}
	}
	invTot := 0
	if tot > 0 {
		invTot = fpow(tot, P-2)
	}
	for i := 1; i <= n; i++ {
		if a[i] == -1 {
			cnt--
		} else {
			leftMiss := sum[a[i]]
			rightMiss := tot - leftMiss
			// expected inversions with -1
			add := (leftMiss*cnt%P + rightMiss*(tot-cnt)%P) % P
			ans = (ans + add*invTot%P) % P
		}
	}
	if ans < 0 {
		ans += P
	}
	fmt.Fprintln(out, ans)
}
