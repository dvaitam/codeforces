package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var m int
		fmt.Fscan(reader, &m)
		top := make([]int, m+1)
		bottom := make([]int, m+1)
		for i := 1; i <= m; i++ {
			fmt.Fscan(reader, &top[i])
		}
		for i := 1; i <= m; i++ {
			fmt.Fscan(reader, &bottom[i])
		}

		prefTop := make([]int, m+1)
		prefBottom := make([]int, m+1)
		for i := 1; i <= m; i++ {
			prefTop[i] = max(prefTop[i-1]+1, top[i]+1)
			prefBottom[i] = max(prefBottom[i-1]+1, bottom[i]+1)
		}
		sufTop := make([]int, m+2)
		sufBottom := make([]int, m+2)
		for i := m; i >= 1; i-- {
			sufTop[i] = max(sufTop[i+1]+1, top[i]+1)
			sufBottom[i] = max(sufBottom[i+1]+1, bottom[i]+1)
		}

		ans := int(1e18)
		for i := 1; i <= m; i++ {
			cur1 := max(prefTop[i-1], sufBottom[i+1])
			cur2 := max(prefBottom[i-1], sufTop[i+1])
			cur := cur1
			if cur2 < cur {
				cur = cur2
			}
			if cur < ans {
				ans = cur
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
