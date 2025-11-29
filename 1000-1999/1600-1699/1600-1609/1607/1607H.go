package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Dish represents the properties of a dish and the calculated constraints
type Dish struct {
	id      int
	a, b, m int
	L, R, S int // L: min remaining fish, R: max remaining fish, S: total remaining weight
}

// toInt parses a byte slice to an integer
func toInt(b []byte) int {
	val := 0
	for _, c := range b {
		val = val*10 + int(c-'0')
	}
	return val
}

func main() {
	// Setup fast I/O
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Read number of test cases
	if !scanner.Scan() {
		return
	}
	t := toInt(scanner.Bytes())

	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			break
		}
		n := toInt(scanner.Bytes())

		dishes := make([]Dish, n)
		for j := 0; j < n; j++ {
			scanner.Scan()
			a := toInt(scanner.Bytes())
			scanner.Scan()
			b := toInt(scanner.Bytes())
			scanner.Scan()
			m := toInt(scanner.Bytes())

			// The total remaining weight for dish j is S = a + b - m.
			// Dishes can only be equal if their S is the same.
			// Let A be the remaining amount of fish.
			// Constraints:
			// 1. 0 <= x <= a  =>  0 <= a - A <= a  =>  0 <= A <= a
			// 2. 0 <= y <= b  =>  0 <= m - (a - A) <= b
			//                 =>  m - b <= a - A <= m
			//                 =>  a - m <= A <= a - m + b
			//                 =>  A >= a - m
			//
			// Also A <= S (since remaining meat >= 0 implies S - A >= 0)
			// So A is in [max(0, a - m), min(a, S)]

			s := a + b - m
			l := a - m
			if l < 0 {
				l = 0
			}
			r := a
			if s < r {
				r = s
			}

			dishes[j] = Dish{id: j, a: a, b: b, m: m, L: l, R: r, S: s}
		}

		// Sort dishes to apply greedy strategy for interval covering
		// Primary key: S (group dishes that can potentially be equal)
		// Secondary key: R (standard greedy interval covering sort)
		sort.Slice(dishes, func(i, j int) bool {
			if dishes[i].S != dishes[j].S {
				return dishes[i].S < dishes[j].S
			}
			return dishes[i].R < dishes[j].R
		})

		ansX := make([]int, n)
		ansY := make([]int, n)
		variety := 0

		// Process each group of dishes with the same S
		for j := 0; j < n; {
			k := j
			// Find the range [j, k) having the same S
			for k < n && dishes[k].S == dishes[j].S {
				k++
			}

			// Greedy interval point cover
			// currentVal stores the chosen 'remaining fish' amount for the current cluster
			currentVal := -1

			for idx := j; idx < k; idx++ {
				d := dishes[idx]
				
				// We need to pick a value in [d.L, d.R].
				// If the last picked value (currentVal) is within this range, we reuse it (minimize variety).
				// Since we sorted by R, and update currentVal greedily to R, we know currentVal <= d.R usually holds
				// if it was picked for a previous interval in the same group.
				// So we strictly check if currentVal >= d.L.
				
				if currentVal == -1 || currentVal < d.L {
					// Cannot reuse, start a new cluster for this group
					currentVal = d.R
					variety++
				}
				
				// Assign eating amounts based on the chosen remaining fish 'currentVal'
				// eaten fish x = initial fish a - remaining fish
				x := d.a - currentVal
				y := d.m - x
				ansX[d.id] = x
				ansY[d.id] = y
			}
			j = k
		}

		fmt.Fprintln(writer, variety)
		for j := 0; j < n; j++ {
			fmt.Fprintln(writer, ansX[j], ansY[j])
		}
	}
}