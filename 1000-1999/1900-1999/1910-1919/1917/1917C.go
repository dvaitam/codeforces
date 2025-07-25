package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// occTime returns the moment (number of prefix operations) when the m-th
// occurrence of value >= j appears in the infinite sequence starting from
// offset s (0-indexed within [0,k)). If it doesn't exist within limits,
// a large number is returned.
func occTime(pos []int, freq, k, s, m int) int {
	if freq == 0 {
		return 1 << 60
	}
	idx := sort.SearchInts(pos, s+1)
	// occurrences left in current cycle after position s
	remain := freq - idx
	if m <= remain {
		return pos[idx+m-1] - s
	}
	t := k - s
	m -= remain
	cycles := (m - 1) / freq
	t += cycles * k
	m -= cycles * freq
	t += pos[m-1]
	return t
}

// bestScore finds the best time t (0 <= t <= L) to perform evaluation
// starting from array a and day offset, returning gained score and t.
func bestScore(a []int, offset, L int, freq []int, occ [][]int, k int) (int, int) {
	type event struct{ t, d int }
	events := make([]event, 0, len(a)*2)
	n := len(a) - 1
	for j := 1; j <= n; j++ {
		need := j - a[j]
		if need < 0 {
			continue
		}
		f := freq[j]
		if f == 0 {
			continue
		}
		tNeed := 0
		if need > 0 {
			tNeed = occTime(occ[j], f, k, offset, need)
		}
		if tNeed > L {
			continue
		}
		tNext := occTime(occ[j], f, k, offset, need+1)
		if tNext > L+1 {
			tNext = L + 1
		}
		events = append(events, event{tNeed, 1}, event{tNext, -1})
	}
	if len(events) == 0 {
		return 0, 0
	}
	sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })
	cur, best, bestT := 0, 0, 0
	for i := 0; i < len(events); {
		t := events[i].t
		for i < len(events) && events[i].t == t {
			cur += events[i].d
			i++
		}
		if cur > best {
			best = cur
			bestT = t
		}
	}
	return best, bestT
}

func solve(n, k int, d int, a []int, v []int) int {
	occ := make([][]int, n+1)
	freq := make([]int, n+1)
	for i, p := range v {
		if p > n {
			p = n
		}
		for j := 1; j <= p; j++ {
			freq[j]++
			occ[j] = append(occ[j], i+1)
		}
	}
	score := 0
	offset := 0
	days := d
	arr := make([]int, n+1)
	copy(arr[1:], a)
	for days > 0 {
		gain, t := bestScore(arr, offset, days-1, freq, occ, k)
		if gain == 0 {
			break
		}
		score += gain
		days -= t + 1
		offset = (offset + t + 1) % k
		for i := 1; i <= n; i++ {
			arr[i] = 0
		}
	}
	return score
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k, d int
		fmt.Fscan(in, &n, &k, &d)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		v := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &v[i])
		}
		fmt.Fprintln(out, solve(n, k, d, a, v))
	}
}
