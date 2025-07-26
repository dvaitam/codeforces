package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const N = 100000
const K = 17 // maximum k needed

var phi [N + 1]int64
var pref [N + 1]int64

// data for each r
type rData struct {
	r      int
	start  []int
	end    []int
	prefix []int64
}

var datas [N + 1]rData

func initPhi() {
	for i := 0; i <= N; i++ {
		phi[i] = int64(i)
	}
	for i := 2; i <= N; i++ {
		if phi[i] == int64(i) {
			for j := i; j <= N; j += i {
				phi[j] -= phi[j] / int64(i)
			}
		}
	}
	for i := 1; i <= N; i++ {
		pref[i] = pref[i-1] + phi[i]
	}
}

func buildData() {
	for r := 1; r <= N; r++ {
		k := 1
		sum := int64(0)
		lim := 2*int(math.Sqrt(float64(r))) + 2
		starts := make([]int, 0, lim)
		ends := make([]int, 0, lim)
		prefix := make([]int64, 0, lim)
		for k <= r {
			q := r / k
			next := r / q
			if next > r {
				next = r
			}
			starts = append(starts, k)
			ends = append(ends, next)
			sum += int64(next-k+1) * pref[q]
			prefix = append(prefix, sum)
			k = next + 1
		}
		datas[r] = rData{r: r, start: starts, end: ends, prefix: prefix}
	}
}

func (d *rData) prefixVal(k int) int64 {
	if k <= 0 {
		return 0
	}
	idx := sort.Search(len(d.start), func(i int) bool { return d.start[i] > k }) - 1
	if idx < 0 {
		return 0
	}
	sum := d.prefix[idx]
	if k < d.end[idx] {
		q := d.r / d.start[idx]
		diff := d.end[idx] - k
		sum -= int64(diff) * pref[q]
	}
	return sum
}

func cost(l, r int) int64 {
	if l > r {
		return 0
	}
	d := &datas[r]
	return d.prefixVal(r) - d.prefixVal(l-1)
}

var dp [K + 1][N + 1]int64

const inf int64 = 1 << 62

func solve(k, l, r, optL, optR int) {
	if l > r {
		return
	}
	mid := (l + r) >> 1
	bestPos := -1
	bestVal := int64(inf)
	maxI := mid - 1
	if optR < maxI {
		maxI = optR
	}
	for i := optL; i <= maxI; i++ {
		val := dp[k-1][i] + cost(i+1, mid)
		if val < bestVal {
			bestVal = val
			bestPos = i
		}
	}
	if bestPos == -1 {
		bestPos = optL
	}
	dp[k][mid] = bestVal
	if l == r {
		return
	}
	solve(k, l, mid-1, optL, bestPos)
	solve(k, mid+1, r, bestPos, optR)
}

func main() {
	initPhi()
	buildData()
	for i := 0; i <= N; i++ {
		for j := 0; j <= K; j++ {
			dp[j][i] = inf
		}
	}
	dp[0][0] = 0
	dp[1][0] = 0
	for n := 1; n <= N; n++ {
		dp[1][n] = cost(1, n)
	}
	for k := 2; k <= K; k++ {
		dp[k][0] = 0
		solve(k, 1, N, 0, N-1)
	}

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		if k > K {
			fmt.Fprintln(out, 0)
			continue
		}
		fmt.Fprintln(out, dp[k][n])
	}
}
