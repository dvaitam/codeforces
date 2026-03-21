package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		if !scanner.Scan() {
			return 0
		}
		res := 0
		for _, b := range scanner.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	t := nextInt()
	for tc := 0; tc < t; tc++ {
		n := nextInt()
		d := nextInt()
		a := make([]int, n)
		vals := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = nextInt()
			vals[i] = a[i]
		}

		sort.Ints(vals)
		m := 0
		for i := 0; i < n; i++ {
			if i == 0 || vals[i] != vals[i-1] {
				vals[m] = vals[i]
				m++
			}
		}
		vals = vals[:m]

		v := make([]int, m+1)
		v[0] = 0
		for i := 0; i < m; i++ {
			v[i+1] = vals[i]
		}

		min_E := make([]int, 4*(m+1))
		max_E := make([]int, 4*(m+1))
		lazy := make([]int, 4*(m+1))
		weight := make([]int64, 4*(m+1))

		var build func(u, l, r int)
		build = func(u, l, r int) {
			lazy[u] = -1
			min_E[u] = 0
			max_E[u] = 0
			if l == r {
				weight[u] = int64(v[l] - v[l-1])
				return
			}
			mid := (l + r) / 2
			build(2*u, l, mid)
			build(2*u+1, mid+1, r)
			weight[u] = weight[2*u] + weight[2*u+1]
		}
		build(1, 1, m)

		var update func(u, l, r, ql, qr, i, val int) int64
		update = func(u, l, r, ql, qr, i, val int) int64 {
			if l > qr || r < ql || min_E[u] > i {
				return 0
			}
			if l >= ql && r <= qr && max_E[u] <= i {
				lazy[u] = val
				min_E[u] = val
				max_E[u] = val
				return weight[u]
			}
			if lazy[u] != -1 {
				lazy[2*u] = lazy[u]
				min_E[2*u] = lazy[u]
				max_E[2*u] = lazy[u]

				lazy[2*u+1] = lazy[u]
				min_E[2*u+1] = lazy[u]
				max_E[2*u+1] = lazy[u]

				lazy[u] = -1
			}
			mid := (l + r) / 2
			var ans int64 = 0
			ans += update(2*u, l, mid, ql, qr, i, val)
			ans += update(2*u+1, mid+1, r, ql, qr, i, val)

			min_E[u] = min_E[2*u]
			if min_E[2*u+1] < min_E[u] {
				min_E[u] = min_E[2*u+1]
			}
			max_E[u] = max_E[2*u]
			if max_E[2*u+1] > max_E[u] {
				max_E[u] = max_E[2*u+1]
			}
			return ans
		}

		var totalAns int64 = 0
		for i := 1; i <= n; i++ {
			target := a[i-1]
			low, high := 1, m
			R := -1
			for low <= high {
				mid := (low + high) / 2
				if v[mid] == target {
					R = mid
					break
				} else if v[mid] < target {
					low = mid + 1
				} else {
					high = mid - 1
				}
			}
			totalAns += update(1, 1, m, 1, R, i, i+d)
		}
		fmt.Fprintln(writer, totalAns)
	}
}
