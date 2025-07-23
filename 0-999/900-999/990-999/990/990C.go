package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	pos := make(map[int]int)
	neg := make(map[int]int)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		bal := 0
		minBal := 0
		for _, ch := range s {
			if ch == '(' {
				bal++
			} else {
				bal--
			}
			if bal < minBal {
				minBal = bal
			}
		}
		if bal >= 0 {
			if minBal >= 0 {
				pos[bal]++
			}
		} else {
			if minBal >= bal {
				neg[-bal]++
			}
		}
	}
	ans := pos[0] * pos[0]
	for k, v := range pos {
		if k == 0 {
			continue
		}
		ans += v * neg[k]
	}
	fmt.Println(ans)
}
