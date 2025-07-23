package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	if len(s) != 6 {
		return
	}

	leftCost := make([]int, 28)
	rightCost := make([]int, 28)
	for i := range leftCost {
		leftCost[i] = 7
		rightCost[i] = 7
	}

	for a := 0; a <= 9; a++ {
		for b := 0; b <= 9; b++ {
			for c := 0; c <= 9; c++ {
				sum := a + b + c
				cost := 0
				if int(s[0]-'0') != a {
					cost++
				}
				if int(s[1]-'0') != b {
					cost++
				}
				if int(s[2]-'0') != c {
					cost++
				}
				if cost < leftCost[sum] {
					leftCost[sum] = cost
				}
			}
		}
	}

	for a := 0; a <= 9; a++ {
		for b := 0; b <= 9; b++ {
			for c := 0; c <= 9; c++ {
				sum := a + b + c
				cost := 0
				if int(s[3]-'0') != a {
					cost++
				}
				if int(s[4]-'0') != b {
					cost++
				}
				if int(s[5]-'0') != c {
					cost++
				}
				if cost < rightCost[sum] {
					rightCost[sum] = cost
				}
			}
		}
	}

	ans := 7
	for sum := 0; sum < 28; sum++ {
		cost := leftCost[sum] + rightCost[sum]
		if cost < ans {
			ans = cost
		}
	}
	fmt.Println(ans)
}
