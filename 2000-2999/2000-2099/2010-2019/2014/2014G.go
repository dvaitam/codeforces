package main

import (
	"bufio"
	"fmt"
	"os"
)

type batch struct {
	day int64
	amt int64
}

const shrinkLimit = 2048

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, k int64
		fmt.Fscan(in, &n, &m, &k)
		days := make([]int64, n)
		amt := make([]int64, n)
		for i := int64(0); i < n; i++ {
			fmt.Fscan(in, &days[i], &amt[i])
		}

		// Deque indices: batches[l:]
		batches := make([]batch, 0, n)
		l := 0
		var total int64
		var ans int64

		// Helper to drop expired milk at the start of day cur.
		expire := func(cur int64) {
			for l < len(batches) && batches[l].day+k <= cur {
				total -= batches[l].amt
				l++
			}
			if l > shrinkLimit {
				batches = append([]batch(nil), batches[l:]...)
				l = 0
			}
		}

		// Helper to consume given amount starting from freshest batches (back).
		consume := func(amount int64) {
			for amount > 0 && len(batches) > l {
				i := len(batches) - 1
				if batches[i].amt <= amount {
					amount -= batches[i].amt
					total -= batches[i].amt
					batches = batches[:i]
				} else {
					batches[i].amt -= amount
					total -= amount
					amount = 0
				}
			}
		}

		if n == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		cur := days[0]
		idx := 0
		for idx < int(n) || l < len(batches) {
			// Remove expired milk for the current day.
			expire(cur)

			// Add new milk arriving today.
			for idx < int(n) && days[idx] == cur {
				batches = append(batches, batch{day: days[idx], amt: amt[idx]})
				total += amt[idx]
				idx++
			}

			// Determine next event day.
			nextDay := int64(1 << 60) // effectively inf
			if idx < int(n) {
				nextDay = days[idx]
			}
			if l < len(batches) {
				expDay := batches[l].day + k
				if expDay < nextDay {
					nextDay = expDay
				}
			}
			if nextDay == int64(1<<60) {
				break
			}

			// Process interval [cur, nextDay)
			delta := nextDay - cur
			for delta > 0 && total > 0 {
				if total >= m {
					full := total / m
					if full > delta {
						full = delta
					}
					ans += full
					consume(full * m)
					delta -= full
				} else {
					// total < m, he drinks everything this day but is not satisfied.
					consume(total)
					delta--
				}
			}

			cur = nextDay
		}

		fmt.Fprintln(out, ans)
	}
}
