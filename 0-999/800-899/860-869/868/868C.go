package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	maxMask := 1 << k
	exists := make([]bool, maxMask)
	for i := 0; i < n; i++ {
		mask := 0
		for j := 0; j < k; j++ {
			var x int
			fmt.Fscan(in, &x)
			if x == 1 {
				mask |= 1 << j
			}
		}
		exists[mask] = true
	}
	masks := make([]int, 0)
	for m := 0; m < maxMask; m++ {
		if exists[m] {
			masks = append(masks, m)
		}
	}
	m := len(masks)
	for subset := 1; subset < (1 << m); subset++ {
		total := 0
		counts := make([]int, k)
		for i := 0; i < m; i++ {
			if subset&(1<<i) != 0 {
				total++
				mask := masks[i]
				for j := 0; j < k; j++ {
					if mask&(1<<j) != 0 {
						counts[j]++
					}
				}
			}
		}
		valid := true
		for j := 0; j < k; j++ {
			if counts[j]*2 > total {
				valid = false
				break
			}
		}
		if valid {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}
