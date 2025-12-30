package main

import (
	"fmt"
)

// ask queries the interactor for the value at index x
func ask(x int) int {
	fmt.Printf("? %d\n", x)
	var y int
	fmt.Scan(&y)
	return y
}

func main() {
	var N int
	fmt.Scan(&N)

	// If N is not divisible by 4, a solution is not guaranteed/possible for this specific problem type
	if N&3 != 0 {
		fmt.Println("! -1")
		return
	}

	l := 1
	r := N/2 + 1

	a := ask(l)
	b := ask(r)

	x := a - b // Difference at the left pointer

	// If the values are already equal (difference is 0), we found the answer at l (1)
	if x == 0 {
		fmt.Println("! 1")
		return
	}

	// Binary Search
	// Invariant: The difference at 'l' (x) and the difference at 'r' have opposite signs.
	for {
		mid := (l + r) >> 1
		
		// Calculate the index opposite to mid: (mid + N/2)
		// We use modular arithmetic adjusted for 1-based indexing
		oppIndex := (mid + N/2 - 1) % N + 1
		
		valMid := ask(mid)
		valOpp := ask(oppIndex)
		z := valMid - valOpp // Difference at mid

		if z == 0 {
			fmt.Printf("! %d\n", mid)
			return
		}

		// Intermediate Value Theorem:
		// If difference at mid (z) has opposite sign to difference at l (x),
		// the zero crossing must be between l and mid.
		if int64(z)*int64(x) < 0 {
			r = mid
			// y = z (conceptually, the value at the right bound updates)
		} else {
			// Otherwise, the zero crossing is between mid and r.
			l = mid
			x = z // Update value at the left bound
		}
	}
}
