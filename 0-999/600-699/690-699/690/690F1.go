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
	deg := make([]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		deg[a]++
		deg[b]++
	}
	var ans int64
	for _, d := range deg {
		ans += int64(d*(d-1)) / 2
	}
	fmt.Println(ans)
}
