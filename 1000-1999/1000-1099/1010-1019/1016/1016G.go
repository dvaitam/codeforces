package main

import (
	"bufio"
	"fmt"
	"os"
)

// Global variables to mirror the logic structure
var (
	n    int
	x, y int64
	a    []int64
	p    []int64 // Stores prime factors of y
	c    []int64 // Array for SOS DP
	cnt  int     // Number of prime factors
)

// Greatest Common Divisor
func gcd(u, v int64) int64 {
	for v != 0 {
		u, v = v, u%v
	}
	return u
}

// spc handles the case where map sizes are large using SOS DP
// Time Complexity: O(Sqrt(Y) + N * Factors + 2^Factors * Factors)
func spc() {
	z := y
	cnt = 0
	p = make([]int64, 0, 20) // Pre-allocate small capacity (max distinct primes for 64-bit int is < 20)

	// 1. Prime factorization of y
	for i := int64(2); i*i <= z; i++ {
		if z%i == 0 {
			for z%i == 0 {
				z /= i
			}
			p = append(p, i)
			cnt++
		}
	}
	if z != 1 {
		p = append(p, z)
		cnt++
	}

	// 2. Populate frequency array c based on divisibility masks
	limit := 1 << cnt
	c = make([]int64, limit)

	for i := 1; i <= n; i++ {
		if a[i]%x != 0 {
			continue
		}
		// We are looking for numbers that can form a valid pair with the first component.
		// num captures the prime factors involved in a[i]/x relative to y/x.
		num := gcd(a[i]/x, y/x)
		msk := 0
		for j := 0; j < cnt; j++ {
			if num%p[j] == 0 {
				msk |= (1 << j)
			}
		}
		c[msk]++
	}

	// 3. SOS DP (Sum Over Subsets)
	// Transforms c[mask] to store the sum of counts for all subsets of mask.
	// This allows O(1) retrieval of "count of numbers disjoint from mask M".
	for i := 0; i < cnt; i++ {
		for j := 0; j < limit; j++ {
			if (j & (1 << i)) == 0 {
				c[j|(1<<i)] += c[j]
			}
		}
	}

	// 4. Calculate Answer
	var ans int64 = 0
	for i := 1; i <= n; i++ {
		if y%a[i] != 0 {
			continue
		}
		// For the second component y/a[i], calculate its prime factor mask
		num := y / a[i]
		msk := 0
		for j := 0; j < cnt; j++ {
			if num%p[j] == 0 {
				msk |= (1 << j)
			}
		}
		// We need gcd(A, B) = 1. If B has factors 'msk', A must NOT have any of those factors.
		// This means A's mask must be a subset of the complement of 'msk'.
		target := (limit - 1) ^ msk
		ans += c[target]
	}
	fmt.Println(ans)
}

func main() {
	// Use bufio for fast I/O
	reader := bufio.NewReader(os.Stdin)

	// Read n, x, y
	fmt.Fscan(reader, &n, &x, &y)

	// Basic validity check
	if y%x != 0 {
		fmt.Println(0)
		return
	}

	// 1-based indexing for convenience to match C++ loop logic
	a = make([]int64, n+1)
	m := make(map[int64]int)
	m2 := make(map[int64]int)

	// Read array a and populate frequency maps
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i]%x == 0 {
			val := gcd(a[i]/x, y/x)
			m[val]++
		}
		if y%a[i] == 0 {
			val := y / a[i]
			m2[val]++
		}
	}

	// Strategy selection based on map size
	if len(m) > 5000 || len(m2) > 5000 {
		spc()
		return
	}

	// Small dataset strategy: Brute force pairs from maps
	var ans int64 = 0
	for k1, v1 := range m {
		for k2, v2 := range m2 {
			if gcd(k1, k2) == 1 {
				ans += int64(v1) * int64(v2)
			}
		}
	}
	fmt.Println(ans)
}
