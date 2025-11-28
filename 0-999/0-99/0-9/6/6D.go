package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, a, b int
   h       []int
   shoot   []int
)

// dfs attempts to use 'balls' shots to reduce all h[1..n] below 0
func dfs(l, balls, last int) bool {
	if balls == 0 {
		for i := 1; i <= n; i++ {
			if h[i] >= 0 {
				return false
			}
		}
		return true
	}
	if l <= n && h[l] < 0 {
		return dfs(l+1, balls, last)
	}
	
	// determine range of possible shots
	start := 2
	end := n - 1
	
	if l > n {
		// All archers from 1 to n are dead, but we have extra balls.
		// We must use them (problem implies exact count or we built up solution).
		// Just shoot at a valid position >= last.
		start = last
		if start < 2 { start = 2 }
	} else {
		// Targeting archer 'l'
		// Valid shots that hurt 'l' are l-1, l, l+1.
		// But l-1 is suboptimal (deals b < a, and l-1 is already dead).
		// So consider l and l+1.
		// If l == n, we can't shoot at n or n+1, so we must shoot at n-1.
		
		if l < n {
			start = l
			end = l + 1
		} else { // l == n
			start = n - 1
			end = n - 1
		}
		
		// Clamp to valid indices [2, n-1]
		if start < 2 { start = 2 }
		if end > n - 1 { end = n - 1 }
		
		// Monotonicity constraint
		if start < last { start = last }
	}
	
	for i := start; i <= end; i++ {
		shoot[balls] = i
		// apply shot
		h[i] -= a
		h[i-1] -= b
		h[i+1] -= b
		if dfs(l, balls-1, i) {
			return true
		}
		// undo shot
		h[i] += a
		h[i-1] += b
		h[i+1] += b
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
		return
	}
	h = make([]int, n+3)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &h[i])
	}
	// allocate shoot positions, size increased to avoid overflow
	shoot = make([]int, 2000)
	// try increasing number of shots
	for ans := 1; ; ans++ {
		if dfs(1, ans, 2) {
			// output result
			fmt.Println(ans)
			for i := ans; i >= 1; i-- {
				fmt.Printf("%d", shoot[i])
				if i > 1 {
					fmt.Printf(" ")
				}
			}
			fmt.Println()
			return
		}
	}
}
