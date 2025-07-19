package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var N int
	if _, err := fmt.Fscan(reader, &N); err != nil {
		return
	}
	pre := 0
	dp := 0
	ans := -1
	for i := 0; i < N; i++ {
		var pos int
		fmt.Fscan(reader, &pos)
		if pos <= pre*2 {
			dp++
		} else {
			dp = 0
		}
		if dp > ans {
			ans = dp
		}
		pre = pos
	}
	fmt.Println(ans + 1)
}
