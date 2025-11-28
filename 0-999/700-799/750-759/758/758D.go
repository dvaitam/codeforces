package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	
	scanner.Scan()
	nStr := scanner.Text()
	n, _ := strconv.ParseUint(nStr, 10, 64)
	
	scanner.Scan()
	k := scanner.Text()
	
	L := len(k)
	// dp[i] stores the minimum value x such that converting x to base n 
	// yields the prefix k[0...i-1]
	dp := make([]uint64, L+1)
	
	// Initialize with max uint64 (infinity representation)
	for i := 1; i <= L; i++ {
		dp[i] = ^uint64(0)
	}
	dp[0] = 0
	
	for i := 1; i <= L; i++ {
		// Try to form the last digit from k[j:i]
		for j := i - 1; j >= 0; j-- {
			// A base-n digit cannot have a leading zero unless it is "0"
			if k[j] == '0' && (i-j) > 1 {
				continue
			}
			
			// Optimization: if substring is too long, it will exceed uint64 or n
			if (i - j) > 20 {
				break
			}
			
			// Parse current digit candidate
			val, err := strconv.ParseUint(k[j:i], 10, 64)
			// If parse fails (overflow) or value is not a valid digit in base n, stop extending
			if err != nil || val >= n {
				break
			}
			
			// If prefix k[0...j] is valid
			if dp[j] != ^uint64(0) {
				prev := dp[j]
				
				// Check for overflow when computing: prev * n + val
				if prev > (^uint64(0))/n {
					continue
				}
				prod := prev * n
				if prod > (^uint64(0)) - val {
					continue
				}
				newVal := prod + val
				
				if newVal < dp[i] {
					dp[i] = newVal
				}
			}
		}
	}
	
	fmt.Println(dp[L])
}