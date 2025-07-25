package main

import (
	"bufio"
	"fmt"
	"os"
)

func nondecreasing(nums []int) bool {
	for i := 1; i < len(nums); i++ {
		if nums[i] < nums[i-1] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		lastDigit := -1
		for i := 0; i < n; i++ {
			if a[i] < 10 {
				lastDigit = i
			}
		}
		possible := false
		for pos := lastDigit + 1; pos <= n; pos++ {
			digits := make([]int, 0, pos*2)
			for i := 0; i < pos; i++ {
				if a[i] >= 10 {
					digits = append(digits, a[i]/10, a[i]%10)
				} else {
					digits = append(digits, a[i])
				}
			}
			if !nondecreasing(digits) {
				continue
			}
			good := true
			for i := pos + 1; i < n; i++ {
				if a[i] < a[i-1] {
					good = false
					break
				}
			}
			if good {
				possible = true
				break
			}
		}
		if possible {
			writer.WriteString("YES\n")
		} else {
			writer.WriteString("NO\n")
		}
	}
}
