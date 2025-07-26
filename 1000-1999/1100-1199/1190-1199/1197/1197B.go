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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// Find position of maximum element
	maxVal := a[0]
	maxIdx := 0
	for i := 1; i < n; i++ {
		if a[i] > maxVal {
			maxVal = a[i]
			maxIdx = i
		}
	}
	// If maximum is at an end, impossible to gather all disks
	if maxIdx == 0 || maxIdx == n-1 {
		fmt.Println("NO")
	} else {
		fmt.Println("YES")
	}
}
