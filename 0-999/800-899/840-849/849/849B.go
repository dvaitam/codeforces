package main

import (
	"bufio"
	"fmt"
	"os"
)

func checkSlope(y []int64, i, j int) bool {
	slopeNumer := y[j-1] - y[i-1]
	slopeDenom := int64(j - i)
	intercepts := make(map[int64]struct{})
	n := len(y)
	for k := 1; k <= n; k++ {
		val := slopeDenom*y[k-1] - slopeNumer*int64(k)
		if _, ok := intercepts[val]; !ok {
			intercepts[val] = struct{}{}
			if len(intercepts) > 2 {
				return false
			}
		}
	}
	return len(intercepts) == 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &y[i])
	}
	if n < 3 {
		fmt.Println("No")
		return
	}
	pairs := [][2]int{{1, 2}, {1, 3}, {2, 3}}
	for _, p := range pairs {
		if checkSlope(y, p[0], p[1]) {
			fmt.Println("Yes")
			return
		}
	}
	fmt.Println("No")
}
