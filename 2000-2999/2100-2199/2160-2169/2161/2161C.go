package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type item struct {
	val int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var x int64
		fmt.Fscan(in, &n, &x)
		arr := make([]int, n)
		items := make([]item, n)
		total := int64(0)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			items[i] = item{val: arr[i], idx: i}
			total += int64(arr[i])
		}
		k := int(total / x)
		if k == 0 {
			fmt.Fprintln(out, 0)
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, arr[i])
			}
			fmt.Fprintln(out)
			continue
		}

		sort.Slice(items, func(i, j int) bool {
			if items[i].val == items[j].val {
				return items[i].idx < items[j].idx
			}
			return items[i].val > items[j].val
		})

		bonus := make([]item, k)
		copy(bonus, items[:k])

		ans := int64(0)
		isBonus := make([]bool, n)
		for i := 0; i < k; i++ {
			ans += int64(bonus[i].val)
			isBonus[bonus[i].idx] = true
		}

		rest := make([]item, 0, n-k)
		for i := 0; i < n; i++ {
			if !isBonus[i] {
				rest = append(rest, item{val: arr[i], idx: i})
			}
		}
		sort.Slice(rest, func(i, j int) bool {
			if rest[i].val == rest[j].val {
				return rest[i].idx < rest[j].idx
			}
			return rest[i].val < rest[j].val
		})

		order := make([]int, 0, n)
		var sum int64
		ri := 0
		for _, b := range bonus {
			need := x - int64(b.val)
			if need < 0 {
				need = 0
			}
			for sum%x < need && ri < len(rest) {
				order = append(order, rest[ri].idx)
				sum += int64(rest[ri].val)
				ri++
			}
			order = append(order, b.idx)
			sum += int64(b.val)
		}
		for ri < len(rest) {
			order = append(order, rest[ri].idx)
			sum += int64(rest[ri].val)
			ri++
		}

		fmt.Fprintln(out, ans)
		for i, idx := range order {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, arr[idx])
		}
		fmt.Fprintln(out)
	}
}
