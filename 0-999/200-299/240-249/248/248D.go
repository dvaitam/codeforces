package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var t int64
	// Read n, t
	// Handle errors gracefully or ignore as per CP style
	fmt.Fscan(reader, &n, &t)
	var s string
	fmt.Fscan(reader, &s)

	// 1-based indexing for convenience matching C++ logic
	// s indices 0..n-1 map to 1..n
	// a array size n+1
	a := make([]int, n+1)
	
	// Precompute 'a'
	// a[i] = a[i-1] + (s[i]=='S') - (s[i]=='H')
	// current balance
	for i := 1; i <= n; i++ {
		val := 0
		if s[i-1] == 'S' {
			val = 1
		} else if s[i-1] == 'H' {
			val = -1
		}
		a[i] = a[i-1] + val
	}

	check := func(x int) bool {
		if int64(a[n]) + int64(x) < 0 {
			return false
		}
		
		to := 0
		// Determine 'to'
		// i from 1 to n
		// a[i-1] + x is balance BEFORE step i?
		// C++: a[i-1]+x < 0
		for i := 1; i <= n; i++ {
			if s[i-1] == 'H' || int64(a[i-1]) + int64(x) < 0 {
				to = i
			}
		}
		
		// Iterate i from 0 to n
		// But wait, loop i from 0 to n-1 represents edges?
		// C++ loop: for(int i=0; i<=n; i++)
		// Checks condition, then updates res.
		// If i=n, it checks condition for ending at n?
		// And updates res for edge n->n+1? (which doesn't exist)
		// Wait, C++ output is based on return inside loop.
		// If loop finishes without return, returns 0.
		// The loop goes up to n.
		
		var res int64 = 0
		for i := 0; i <= n; i++ {
			// Cost if we stop at i and do excursion to 'to'
			// Dist from i to 'to' is max(to-i, 0)
			// Round trip cost 2 * dist
			// Total cost res + 2*dist
			// Note: res includes cost to reach i.
			// Is cost to reach 0 is 0?
			// Loop i=0: check cost if stop at 0 (with excursion to to).
			// res is 0 initially.
			
			dist := to - i
			if dist < 0 {
				dist = 0
			}
			if res + int64(dist)*2 <= t {
				return true
			}
			
			// Update res for next step (i -> i+1)
			// If i == n, we shouldn't update or it doesn't matter
			if i == n {
				break
			}
			
			cost := int64(1)
			// If deficit at i (after processing i-th element), we traverse edge i->i+1 3 times (add 2)
			// C++: res += (a[i]+x < 0)*2 + 1
			if int64(a[i]) + int64(x) < 0 {
				cost += 2
			}
			res += cost
		}
		return false
	}

	// Binary search
	l, r := 0, n+1
	for l < r {
		mid := (l + r) >> 1
		if check(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}

	if l <= n {
		fmt.Println(l)
	} else {
		fmt.Println("-1")
	}
}