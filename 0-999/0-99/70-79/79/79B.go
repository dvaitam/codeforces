package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int64
	var k, t int
	if _, err := fmt.Fscan(reader, &n, &m, &k, &t); err != nil {
		return
	}

	waste := make([]int64, k)
	wasteSet := make(map[int64]bool)
	for i := 0; i < k; i++ {
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		pos := (a-1)*m + b
		waste[i] = pos
		wasteSet[pos] = true
	}
	
	sort.Slice(waste, func(i, j int) bool { return waste[i] < waste[j] })
	
	for qi := 0; qi < t; qi++ {
		var x, y int64
		fmt.Fscan(reader, &x, &y)
		pos := (x-1)*m + y
		if wasteSet[pos] {
			fmt.Println("Waste")
		} else {
			// count waste cells before this position
			cnt := sort.Search(len(waste), func(i int) bool { return waste[i] >= pos })
			r := (pos - int64(cnt)) % 3
			if r == 1 {
				fmt.Println("Carrots")
			} else if r == 2 {
				fmt.Println("Kiwis")
			} else {
				fmt.Println("Grapes")
			}
		}
	}
}
