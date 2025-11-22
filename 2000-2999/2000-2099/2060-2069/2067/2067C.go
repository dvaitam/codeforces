package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxOps = 8

var adds []int64
var addSums [maxOps + 1][]int64

func buildAdds() {
	adds = make([]int64, 10) // lengths 1..10
	x := int64(0)
	for i := 0; i < 10; i++ {
		x = x*10 + 9
		adds[i] = x
	}
}

func genSums() {
	var dfs func(idx, rem, k int, cur int64)
	dfs = func(idx, rem, k int, cur int64) {
		if idx == len(adds)-1 {
			addSums[k] = append(addSums[k], cur+adds[idx]*int64(rem))
			return
		}
		for c := 0; c <= rem; c++ {
			dfs(idx+1, rem-c, k, cur+adds[idx]*int64(c))
		}
	}
	addSums[0] = []int64{0}
	for k := 1; k <= maxOps; k++ {
		dfs(0, k, k, 0)
	}
}

func hasSeven(x int64) bool {
	for x > 0 {
		if x%10 == 7 {
			return true
		}
		x /= 10
	}
	return false
}

func main() {
	buildAdds()
	genSums()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int64
		fmt.Fscan(in, &n)
		if hasSeven(n) {
			fmt.Fprintln(out, 0)
			continue
		}
		ans := maxOps + 1
		for k := 1; k <= maxOps; k++ {
			found := false
			for _, add := range addSums[k] {
				if hasSeven(n + add) {
					found = true
					break
				}
			}
			if found {
				ans = k
				break
			}
		}
		// With the tested bound maxOps, the answer is always found.
		fmt.Fprintln(out, ans)
	}
}
