package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func checkMin(targetMin int64, k int64, c []int64) bool {
	needed := int64(0)
	for _, wealth := range c {
		if wealth < targetMin {
			needed += targetMin - wealth
		}
	}
	return needed <= k
}

func checkMax(targetMax int64, k int64, c []int64) bool {
	removed := int64(0)
	for _, wealth := range c {
		if wealth > targetMax {
			removed += wealth - targetMax
		}
	}
	return removed <= k
}

func solve() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line1 := scanner.Text()
	parts := strings.Fields(line1)
	n, _ := strconv.Atoi(parts[0])
	k, _ := strconv.ParseInt(parts[1], 10, 64)
	scanner.Scan()
	line2 := scanner.Text()
	wealthStr := strings.Fields(line2)
	c := make([]int64, n)
	for i, s := range wealthStr {
		c[i], _ = strconv.ParseInt(s, 10, 64)
	}
	if n <= 1 {
		fmt.Println("0")
		return
	}
	totalSum := int64(0)
	for _, v := range c {
		totalSum += v
	}
	n64 := int64(n)
	minVal := c[0]
	maxVal := c[0]
	for _, v := range c {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	maxMinPossible := totalSum / n64
	minMaxPossible := (totalSum + n64 - 1) / n64
	low := minVal
	high := maxMinPossible
	finalMinWealth := minVal
	for low <= high {
		mid := low + (high-low)/2
		if checkMin(mid, k, c) {
			finalMinWealth = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	low = minMaxPossible
	high = maxVal
	finalMaxWealth := maxVal
	for low <= high {
		mid := low + (high-low)/2
		if checkMax(mid, k, c) {
			finalMaxWealth = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	fmt.Println(finalMaxWealth - finalMinWealth)
}

func main() {
	solve()
}
