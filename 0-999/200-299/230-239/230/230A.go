package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s, n int
	fmt.Fscan(reader, &s, &n)
	dragons := make([]struct{ x, y int }, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &dragons[i].x, &dragons[i].y)
	}
	sort.Slice(dragons, func(i, j int) bool {
		return dragons[i].x < dragons[j].x
	})
	for _, d := range dragons {
		if s <= d.x {
			fmt.Println("NO")
			return
		}
		s += d.y
	}
	fmt.Println("YES")
}
