package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Task struct {
	k int64
	b int64
	a int64
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
		tasks := make([]Task, n)
		var base int64
		for i := 0; i < n; i++ {
			var k, b, a int64
			fmt.Fscan(reader, &k, &b, &a)
			tasks[i] = Task{k: k, b: b, a: a}
			base += a
		}
		sort.Slice(tasks, func(i, j int) bool {
			if tasks[i].k == tasks[j].k {
				return tasks[i].b > tasks[j].b
			}
			return tasks[i].k > tasks[j].k
		})
		const inf int64 = -1 << 60
		dp := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			dp[i] = inf
		}
		for idx, t := range tasks {
			c := t.b - t.a
			if c <= 0 {
				continue
			}
			d := (c - 1) / t.k
			if int(d) > n {
				d = int64(n)
			}
			limit := idx + 1
			if limit > n {
				limit = n
			}
			if limit > int(d) {
				limit = int(d)
			}
			for m := limit; m >= 1; m-- {
				if dp[m-1] == inf {
					continue
				}
				cand := dp[m-1] + c - t.k*int64(m)
				if cand > dp[m] {
					dp[m] = cand
				}
			}
		}
		var best int64
		for i := 0; i <= n; i++ {
			if dp[i] > best {
				best = dp[i]
			}
		}
		fmt.Fprintln(writer, base+best)
	}
}
