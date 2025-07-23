package main

import (
	"bufio"
	"fmt"
	"os"
)

// countZeros returns the number of trailing zeros in n!
func countZeros(n int) int {
	count := 0
	for n > 0 {
		n /= 5
		count += n
	}
	return count
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	low, high := 0, 5*m
	for low < high {
		mid := (low + high) / 2
		if countZeros(mid) < m {
			low = mid + 1
		} else {
			high = mid
		}
	}
	if countZeros(low) != m {
		fmt.Println(0)
		return
	}
	fmt.Println(5)
	for i := 0; i < 5; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(low + i)
	}
	fmt.Println()
}
