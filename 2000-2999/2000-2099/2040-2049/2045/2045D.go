package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	llinf = int64(1) << 40
	segM  = 1 << 20
)

var seg [segM*2 + 10][2]int64

func mer(a, b [2]int64) [2]int64 {
	y := a[1]
	if a[0]+b[1] < y {
		y = a[0] + b[1]
	}
	return [2]int64{a[0] + b[0], y}
}

func ge(l, r, lo, hi, i int) [2]int64 {
	if lo >= l && hi <= r {
		return seg[i]
	}
	if lo >= r || hi <= l {
		return [2]int64{0, llinf * 1000010}
	}
	mid := (lo + hi) / 2
	return mer(ge(l, r, lo, mid, i*2), ge(l, r, mid, hi, i*2+1))
}

func upd(i int, x int64) {
	i += segM
	seg[i][0] += x
	seg[i][1] += x
	for i /= 2; i > 0; i /= 2 {
		seg[i] = mer(seg[i*2], seg[i*2+1])
	}
}

func getInd(vas []int64, x int64) int {
	return sort.Search(len(vas), func(i int) bool { return vas[i] >= x })
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var d, sw, fl, wa int64
	fmt.Fscan(in, &n, &d, &sw, &fl, &wa)

	if fl < sw {
		fmt.Fprintln(out, int64(n-1)*fl)
		return
	}

	pr := make([]int64, n+1)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		pr[i+1] = pr[i] + x - d
	}

	var vas []int64
	for i := 0; i <= n; i++ {
		vas = append(vas, pr[i], pr[i]+d, pr[i]+2*d)
	}
	vas = append(vas, -llinf)
	sort.Slice(vas, func(i, j int) bool { return vas[i] < vas[j] })
	j := 1
	for k := 1; k < len(vas); k++ {
		if vas[k] != vas[k-1] {
			vas[j] = vas[k]
			j++
		}
	}
	vas = vas[:j]

	upd(0, llinf)
	idx0 := getInd(vas, 0)
	upd(idx0, -llinf)
	upd(idx0+1, llinf)

	dp := make([]int64, n+1)

	for i := 1; i < n; i++ {
		dp[i] = ge(0, getInd(vas, pr[i]+d), 0, segM, 1)[1] + int64(i-1)*sw + fl
		swimtohere := ge(0, getInd(vas, pr[i]+d)+1, 0, segM, 1)[1] + int64(i-1)*sw
		flytohere := ge(0, getInd(vas, pr[i]+2*d), 0, segM, 1)[1] + int64(i-2)*sw + fl

		p2 := getInd(vas, pr[i]+d)
		va2 := min64(swimtohere, flytohere) - int64(i-1)*sw
		cu2 := ge(0, p2+1, 0, segM, 1)[0]
		if va2 < cu2 {
			upd(p2, va2-cu2)
			upd(p2+1, cu2-va2)
		}

		pl := getInd(vas, pr[i])
		upd(pl+1, 2*wa)
		va := dp[i] - int64(i)*sw
		cu := ge(0, pl+1, 0, segM, 1)[0]
		if va < cu {
			upd(pl, va-cu)
			upd(pl+1, cu-va)
		}
	}

	swimtoend := ge(0, getInd(vas, pr[n]+d)+1, 0, segM, 1)[1] + int64(n-1)*sw
	flytoend := ge(0, getInd(vas, pr[n]+2*d), 0, segM, 1)[1] + int64(n-2)*sw + fl

	ans := swimtoend
	if flytoend < ans {
		ans = flytoend
	}
	fmt.Fprintln(out, ans)
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
