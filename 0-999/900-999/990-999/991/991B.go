package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	grades := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grades[i])
		sum += grades[i]
	}

	target := 45 * n
	cur := sum * 10
	if cur >= target {
		fmt.Println(0)
		return
	}

	sort.Ints(grades)
	cnt := 0
	for _, g := range grades {
		diff := 5 - g
		cur += diff * 10
		cnt++
		if cur >= target {
			fmt.Println(cnt)
			return
		}
	}
}
