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
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}

	max1Val := -1
	max2Val := -1
	max1Idx := -1
	inc := make([]int, n)
	isRec := make([]bool, n)
	records := 0

	for i := 0; i < n; i++ {
		v := p[i]
		if v > max1Val {
			records++
			isRec[i] = true
			max2Val = max1Val
			max1Val = v
			max1Idx = i
		} else if v > max2Val {
			inc[max1Idx]++
			max2Val = v
		}
	}

	bestRecords := -1
	bestVal := int(^uint(0) >> 1)
	for i := 0; i < n; i++ {
		r := records
		if isRec[i] {
			r--
		}
		r += inc[i]
		if r > bestRecords || (r == bestRecords && p[i] < bestVal) {
			bestRecords = r
			bestVal = p[i]
		}
	}

	fmt.Println(bestVal)
}
