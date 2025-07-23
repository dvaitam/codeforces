package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	cnt := make([]int, 10)
	for _, ch := range s {
		cnt[ch-'0']++
	}

	digits := []int{}
	for d := 0; d < 10; d++ {
		if cnt[d] > 0 {
			digits = append(digits, d)
		}
	}

	maxLen := len(s)
	fact := make([]int64, maxLen+1)
	fact[0] = 1
	for i := 1; i <= maxLen; i++ {
		fact[i] = fact[i-1] * int64(i)
	}

	cur := make([]int, 10)
	var ans int64
	var dfs func(int, int)
	dfs = func(pos int, length int) {
		if pos == len(digits) {
			L := length
			total := fact[L]
			for _, d := range digits {
				total /= fact[cur[d]]
			}
			if cur[0] > 0 {
				t := fact[L-1]
				t /= fact[cur[0]-1]
				for _, d := range digits {
					if d == 0 {
						continue
					}
					t /= fact[cur[d]]
				}
				total -= t
			}
			ans += total
			return
		}
		d := digits[pos]
		for c := 1; c <= cnt[d]; c++ {
			cur[d] = c
			dfs(pos+1, length+c)
		}
		cur[d] = 0
	}

	dfs(0, 0)
	fmt.Println(ans)
}
