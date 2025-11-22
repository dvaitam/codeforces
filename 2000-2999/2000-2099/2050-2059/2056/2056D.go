package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type point struct {
	u int64
	v int
}

type bit struct {
	n    int
	data []int
}

func newBit(n int) *bit {
	return &bit{n: n, data: make([]int, n+2)}
}

func (b *bit) add(idx, delta int) {
	for idx <= b.n {
		b.data[idx] += delta
		idx += idx & -idx
	}
}

func (b *bit) query(idx int) int {
	res := 0
	for idx > 0 {
		res += b.data[idx]
		idx -= idx & -idx
	}
	return res
}

func countPairs(points []point) int64 {
	if len(points) == 0 {
		return 0
	}

	// coordinate compression on v
	allV := make([]int64, len(points))
	for i, p := range points {
		allV[i] = int64(p.v)
	}
	sort.Slice(allV, func(i, j int) bool { return allV[i] < allV[j] })
	allV = unique64(allV)
	for i := range points {
		points[i].v = sort.Search(len(allV), func(j int) bool { return allV[j] >= int64(points[i].v) }) + 1
	}

	a := make([]point, len(points))
	copy(a, points)
	tmp := make([]point, len(points))
	b := newBit(len(allV) + 2)

	var cdq func(l, r int) int64
	cdq = func(l, r int) int64 {
		if l >= r {
			return 0
		}
		mid := (l + r) >> 1
		ans := cdq(l, mid) + cdq(mid+1, r)

		i := l
		used := make([]int, 0)
		for j := mid + 1; j <= r; j++ {
			for i <= mid && a[i].u < a[j].u {
				b.add(a[i].v, 1)
				used = append(used, a[i].v)
				i++
			}
			ans += int64(b.query(a[j].v - 1))
		}
		for _, v := range used {
			b.add(v, -1)
		}

		// merge by u
		p, q := l, mid+1
		k := l
		for p <= mid && q <= r {
			if a[p].u <= a[q].u {
				tmp[k] = a[p]
				p++
			} else {
				tmp[k] = a[q]
				q++
			}
			k++
		}
		for p <= mid {
			tmp[k] = a[p]
			p++
			k++
		}
		for q <= r {
			tmp[k] = a[q]
			q++
			k++
		}
		for i := l; i <= r; i++ {
			a[i] = tmp[i]
		}
		return ans
	}

	return cdq(0, len(a)-1)
}

func unique64(arr []int64) []int64 {
	if len(arr) == 0 {
		return arr
	}
	res := arr[:1]
	for _, v := range arr[1:] {
		if v != res[len(res)-1] {
			res = append(res, v)
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		// prefix counts per value
		pref := make([][]int, 11)
		for i := 0; i <= 10; i++ {
			pref[i] = make([]int, n+1)
		}
		for i := 0; i < n; i++ {
			for v := 1; v <= 10; v++ {
				pref[v][i+1] = pref[v][i]
			}
			pref[a[i]][i+1]++
		}

		// cumulative sum over values
		cum := make([][]int, 11)
		for v := 0; v <= 10; v++ {
			cum[v] = make([]int, n+1)
		}
		for v := 1; v <= 10; v++ {
			for i := 0; i <= n; i++ {
				cum[v][i] = cum[v-1][i] + pref[v][i]
			}
		}

		// count odd subarrays (always good)
		evenPref := n/2 + 1
		oddPref := (n + 1) / 2
		ans := int64(evenPref * oddPref)

		for val := 1; val <= 10; val++ {
			pointsEven := make([]point, 0, (n+2)/2)
			pointsOdd := make([]point, 0, (n+2)/2)
			for i := 0; i <= n; i++ {
				less := cum[val-1][i]
				zero := pref[val][i]
				u := int64(i - 2*less)
				vCoord := int64(2*less + 2*zero - i)
				pt := point{u: u, v: int(vCoord)}
				if i%2 == 0 {
					pointsEven = append(pointsEven, pt)
				} else {
					pointsOdd = append(pointsOdd, pt)
				}
			}
			ans += countPairs(pointsEven)
			ans += countPairs(pointsOdd)
		}

		fmt.Fprintln(out, ans)
	}
}
