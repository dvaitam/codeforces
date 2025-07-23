package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	xCount := make(map[int]int)
	yCount := make(map[int]int)
	pairCount := make(map[[2]int]int)

	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		xCount[x]++
		yCount[y]++
		pairCount[[2]int{x, y}]++
	}

	var ans int64
	for _, c := range xCount {
		ans += int64(c) * int64(c-1) / 2
	}
	for _, c := range yCount {
		ans += int64(c) * int64(c-1) / 2
	}
	for _, c := range pairCount {
		ans -= int64(c) * int64(c-1) / 2
	}

	fmt.Println(ans)
}
