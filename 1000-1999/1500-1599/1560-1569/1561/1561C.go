package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cave struct {
	need int
	k    int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		caves := make([]Cave, n)
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(in, &k)
			maxNeed := 0
			for j := 0; j < k; j++ {
				var a int
				fmt.Fscan(in, &a)
				if val := a - j; val > maxNeed {
					maxNeed = val
				}
			}
			caves[i] = Cave{need: maxNeed + 1, k: k}
		}

		sort.Slice(caves, func(i, j int) bool {
			if caves[i].need == caves[j].need {
				return caves[i].k < caves[j].k
			}
			return caves[i].need < caves[j].need
		})

		start := caves[0].need
		cur := start
		cur += caves[0].k
		for i := 1; i < n; i++ {
			if cur < caves[i].need {
				start += caves[i].need - cur
				cur = caves[i].need
			}
			cur += caves[i].k
		}
		fmt.Fprintln(out, start)
	}
}
