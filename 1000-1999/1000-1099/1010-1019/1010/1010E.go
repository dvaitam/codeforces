package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	MaxN   = 200005
	BITLen = 100005
	Inf    = 1000000000
)

var (
	t          [BITLen]int
	a, b, c    [MaxN]int
	an         [MaxN]int
	n, m, kVal int
	x, X       int
	y, Y       int
	z, Z       int
)

type nd struct {
	u, v, w, id int
}

// add updates the Fenwick tree with prefix minimums
func add(idx, val int) {
	for ; idx < BITLen; idx += idx & -idx {
		if val < t[idx] {
			t[idx] = val
		}
	}
}

// qr queries the Fenwick tree for the prefix minimum
func qr(idx int) int {
	res := Inf
	for ; idx > 0; idx -= idx & -idx {
		if t[idx] < res {
			res = t[idx]
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Use Fscan for relatively fast parsing of large integer sets
	if _, err := fmt.Fscan(reader, &n, &m, &kVal); err != nil {
		return
	}
	// Re-reading to match original C++ logic which overwrites initial variables
	if _, err := fmt.Fscan(reader, &n, &m, &kVal); err != nil {
		return
	}

	x, y, z = Inf, Inf, Inf
	X, Y, Z = -Inf, -Inf, -Inf

	// Calculate the 3D bounding box of the initial points
	for i := 0; i < n; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		if u < x { x = u }
		if u > X { X = u }
		if v < y { y = v }
		if v > Y { Y = v }
		if w < z { z = w }
		if w > Z { Z = w }
	}

	// Read 'forbidden' points (M) and check if any are inside the box
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &a[i], &b[i], &c[i])
		if a[i] >= x && a[i] <= X && b[i] >= y && b[i] <= Y && c[i] >= z && c[i] <= Z {
			fmt.Fprintln(writer, "INCORRECT")
			return
		}
	}

	fmt.Fprintln(writer, "CORRECT")

	// Read query points (K) and check if they are already inside the box (OPEN)
	for i := 1; i <= kVal; i++ {
		fmt.Fscan(reader, &a[i+m], &b[i+m], &c[i+m])
		if a[i+m] >= x && a[i+m] <= X && b[i+m] >= y && b[i+m] <= Y && c[i+m] >= z && c[i+m] <= Z {
			an[i] = 1 // Status: OPEN
		}
	}

	// 3D Dominance check using 8 octant transformations
	// This determines if a query point is "blocked" by a forbidden point
	for s0 := 0; s0 <= 1; s0++ {
		for s1 := 0; s1 <= 1; s1++ {
			for s2 := 0; s2 <= 1; s2++ {
				var d []nd
				for i := 1; i <= m+kVal; i++ {
					// Check if point is outside the bounding box in the current direction
					cond0 := (s0 == 1 && X >= a[i]) || (s0 == 0 && a[i] >= x)
					cond1 := (s1 == 1 && Y >= b[i]) || (s1 == 0 && b[i] >= y)
					cond2 := (s2 == 1 && Z >= c[i]) || (s2 == 0 && c[i] >= z)

					if cond0 && cond1 && cond2 {
						// i <= m: Forbidden point; i > m: Query point that isn't already OPEN
						if i <= m || (i > m && an[i-m] == 0) {
							uVal, vVal, wVal := 0, 0, 0
							if s0 == 1 { uVal = X - a[i] } else { uVal = a[i] - x }
							if s1 == 1 { vVal = Y - b[i] } else { vVal = b[i] - y }
							if s2 == 1 { wVal = Z - c[i] } else { wVal = c[i] - z }
							d = append(d, nd{uVal, vVal, wVal, i})
						}
					}
				}

				// Sort points by transformed first coordinate (u)
				sort.Slice(d, func(i, j int) bool {
					if d[i].u != d[j].u {
						return d[i].u < d[j].u
					}
					return d[i].id < d[j].id
				})

				// Reset Fenwick Tree for each transformation pass
				for i := 0; i < BITLen; i++ {
					t[i] = Inf
				}

				// BIT processing for 3D partial order
				for i := 0; i < len(d); i++ {
					if d[i].id <= m {
						// Add forbidden point to BIT
						add(d[i].v+1, d[i].w)
					} else {
						// Check if query point is dominated by any forbidden point
						if qr(d[i].v+1) <= d[i].w {
							an[d[i].id-m] = 2 // Status: CLOSED
						}
					}
				}
			}
		}
	}

	// Final Output
	for i := 1; i <= kVal; i++ {
		switch an[i] {
		case 1:
			fmt.Fprintln(writer, "OPEN")
		case 2:
			fmt.Fprintln(writer, "CLOSED")
		default:
			fmt.Fprintln(writer, "UNKNOWN")
		}
	}
}
