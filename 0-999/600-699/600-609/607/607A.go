package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	type Beacon struct{ x, p int }
	bs := make([]Beacon, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &bs[i].x, &bs[i].p)
	}
	sort.Slice(bs, func(i, j int) bool { return bs[i].x < bs[j].x })

	pos := make([]int, n)
	pow := make([]int, n)
	for i := 0; i < n; i++ {
		pos[i] = bs[i].x
		pow[i] = bs[i].p
	}

	dp := make([]int, n)
	for i := 0; i < n; i++ {
		left := pos[i] - pow[i]
		j := sort.Search(i, func(k int) bool { return pos[k] >= left }) - 1
		if j >= 0 {
			dp[i] = dp[j] + (i - j - 1)
		} else {
			dp[i] = i
		}
	}

	ans := n
	for i := 0; i < n; i++ {
		destroyed := dp[i] + (n - i - 1)
		if destroyed < ans {
			ans = destroyed
		}
	}
	if n < ans {
		ans = n
	}
	fmt.Println(ans)
}
