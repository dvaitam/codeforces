package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		costs := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &costs[i])
		}

		// initialize using first two segments
		minOdd := costs[0]
		minEven := costs[1]
		sumOdd := costs[0]
		sumEven := costs[1]
		cntOdd := 1
		cntEven := 1
		ans := (int64(n)-int64(cntOdd))*minOdd + (int64(n)-int64(cntEven))*minEven + sumOdd + sumEven

		for i := 2; i < n; i++ {
			v := costs[i]
			if i%2 == 0 {
				sumOdd += v
				cntOdd++
				if v < minOdd {
					minOdd = v
				}
			} else {
				sumEven += v
				cntEven++
				if v < minEven {
					minEven = v
				}
			}
			current := (int64(n)-int64(cntOdd))*minOdd + (int64(n)-int64(cntEven))*minEven + sumOdd + sumEven
			if current < ans {
				ans = current
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
