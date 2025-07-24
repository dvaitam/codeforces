package main

import (
	"bufio"
	"fmt"
	"os"
)

func encode(card []int) int64 {
	var x int64
	for _, v := range card {
		x = x*3 + int64(v)
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	cards := make([][]int, n)
	mp := make(map[int64]int)
	for i := 0; i < n; i++ {
		cards[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &cards[i][j])
		}
		mp[encode(cards[i])] = i
	}

	thr := make([]int, k)
	var result int64
	for i := 0; i < n; i++ {
		cnt := 0
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			for t := 0; t < k; t++ {
				if cards[i][t] == cards[j][t] {
					thr[t] = cards[i][t]
				} else {
					thr[t] = 3 - cards[i][t] - cards[j][t]
				}
			}
			val := int64(0)
			for t := 0; t < k; t++ {
				val = val*3 + int64(thr[t])
			}
			if idx, ok := mp[val]; ok && idx != i && idx != j && j < idx {
				cnt++
			}
		}
		result += int64(cnt * (cnt - 1) / 2)
	}
	fmt.Fprintln(out, result)
}
