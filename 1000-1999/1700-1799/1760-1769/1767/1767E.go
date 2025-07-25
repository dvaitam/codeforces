package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &colors[i])
		colors[i]--
	}

	cost := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &cost[i])
	}

	// build adjacency matrix
	adj := make([]uint64, m)
	for i := 0; i < n-1; i++ {
		a := colors[i]
		b := colors[i+1]
		adj[a] |= 1 << uint(b)
		adj[b] |= 1 << uint(a)
	}

	totalCost := int64(0)
	for i := 0; i < m; i++ {
		totalCost += cost[i]
	}

	forced := make([]bool, m)
	forced[colors[0]] = true
	forced[colors[n-1]] = true
	for i := 0; i < m; i++ {
		if (adj[i]>>uint(i))&1 != 0 {
			forced[i] = true
		}
	}

	// remove forced vertices from graph
	for i := 0; i < m; i++ {
		if forced[i] {
			adj[i] = 0
		}
	}
	for i := 0; i < m; i++ {
		if adj[i] == 0 {
			continue
		}
		mask := uint64(0)
		for j := 0; j < m; j++ {
			if forced[j] {
				mask |= 1 << uint(j)
			}
		}
		adj[i] &^= mask
	}

	idx := make([]int, 0)
	for i := 0; i < m; i++ {
		if !forced[i] {
			idx = append(idx, i)
		}
	}
	r := len(idx)
	if r == 0 {
		fmt.Fprintln(writer, totalCost)
		return
	}

	weights := make([]int64, r)
	for i := 0; i < r; i++ {
		weights[i] = cost[idx[i]]
	}

	matrix := make([]uint64, r)
	for i := 0; i < r; i++ {
		mask := uint64(0)
		for j := 0; j < r; j++ {
			if i != j {
				if (adj[idx[i]]>>uint(idx[j]))&1 != 0 {
					mask |= 1 << uint(j)
				}
			}
		}
		matrix[i] = mask
	}

	n1 := r / 2
	n2 := r - n1

	adj1 := make([]uint64, n1)
	cross1 := make([]uint64, n1)
	w1 := make([]int64, n1)
	for i := 0; i < n1; i++ {
		w1[i] = weights[i]
		var mask1 uint64
		for j := 0; j < n1; j++ {
			if (matrix[i]>>uint(j))&1 != 0 {
				mask1 |= 1 << uint(j)
			}
		}
		adj1[i] = mask1
		var mask2 uint64
		for j := 0; j < n2; j++ {
			if (matrix[i]>>uint(n1+j))&1 != 0 {
				mask2 |= 1 << uint(j)
			}
		}
		cross1[i] = mask2
	}

	adj2 := make([]uint64, n2)
	w2 := make([]int64, n2)
	for i := 0; i < n2; i++ {
		w2[i] = weights[n1+i]
		var mask uint64
		for j := 0; j < n2; j++ {
			if (matrix[n1+i]>>uint(n1+j))&1 != 0 {
				mask |= 1 << uint(j)
			}
		}
		adj2[i] = mask
	}

	size1 := 1 << uint(n2)
	best := make([]int64, size1)
	allMaskB := uint64(size1 - 1)

	for mask := 0; mask < (1 << uint(n1)); mask++ {
		valid := true
		var weight int64
		var forbid uint64
		for i := 0; i < n1 && valid; i++ {
			if (mask>>uint(i))&1 != 0 {
				if (uint64(mask) & adj1[i]) != 0 {
					valid = false
					break
				}
				weight += w1[i]
				forbid |= cross1[i]
			}
		}
		if valid {
			allowed := allMaskB &^ forbid
			if weight > best[allowed] {
				best[allowed] = weight
			}
		}
	}

	for i := 0; i < n2; i++ {
		for mask := 0; mask < size1; mask++ {
			if mask&(1<<i) != 0 {
				if best[mask^(1<<i)] < best[mask] {
					best[mask^(1<<i)] = best[mask]
				}
			}
		}
	}

	var mwis int64
	for mask := 0; mask < (1 << uint(n2)); mask++ {
		valid := true
		var weight int64
		for i := 0; i < n2 && valid; i++ {
			if (mask>>uint(i))&1 != 0 {
				if (uint64(mask) & adj2[i]) != 0 {
					valid = false
					break
				}
				weight += w2[i]
			}
		}
		if valid {
			if weight+best[mask] > mwis {
				mwis = weight + best[mask]
			}
		}
	}

	result := totalCost - mwis
	fmt.Fprintln(writer, result)
}
