package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	if n == 1 {
		if a[0] >= 0 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, -1)
		}
		return
	}

	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}

	target := prefix[1:n] // S_1 .. S_{N-1}
	total := prefix[n]

	minVal, maxVal := target[0], target[0]
	for _, v := range target {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	if minVal < 0 || maxVal > total {
		fmt.Fprintln(out, -1)
		return
	}

	vals := make([]int64, len(target))
	copy(vals, target)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	k := 0
	for i := 0; i < len(vals); i++ {
		if i == 0 || vals[i] != vals[i-1] {
			vals[k] = vals[i]
			k++
		}
	}
	vals = vals[:k]

	bit := make([]int, k+2)
	update := func(idx int) {
		for idx < len(bit) {
			bit[idx]++
			idx += idx & -idx
		}
	}
	query := func(idx int) int {
		sum := 0
		for idx > 0 {
			sum += bit[idx]
			idx -= idx & -idx
		}
		return sum
	}
	rank := func(val int64) int {
		pos := sort.Search(len(vals), func(i int) bool { return vals[i] >= val })
		return pos + 1
	}

	var ans int64
	for i, v := range target {
		r := rank(v)
		lessOrEqual := query(r)
		ans += int64(i - lessOrEqual)
		update(r)
	}

	fmt.Fprintln(out, ans)
}
