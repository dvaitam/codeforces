package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Use buffered I/O for speed
	reader := bufio.NewReader(os.Stdin)

	// Read input strings
	t, _ := reader.ReadString('\n')
	t = strings.TrimSpace(t)

	sBegin, _ := reader.ReadString('\n')
	sBegin = strings.TrimSpace(sBegin)

	sEnd, _ := reader.ReadString('\n')
	sEnd = strings.TrimSpace(sEnd)

	n := len(t)
	nb := len(sBegin)
	ne := len(sEnd)

	// Basic constraint check
	if nb > n || ne > n {
		fmt.Println(0)
		return
	}

	// Rolling Hash Constants
	// Hash 1: Explicit Modulo 10^9 + 7
	mod1 := uint64(1000000007)
	base1 := uint64(131)
	// Hash 2: Implicit Modulo 2^64 (uint64 overflow)
	base2 := uint64(13331)

	// Precompute hash arrays
	h1 := make([]uint64, n+1)
	p1 := make([]uint64, n+1)
	h2 := make([]uint64, n+1)
	p2 := make([]uint64, n+1)

	p1[0] = 1
	p2[0] = 1

	for i := 0; i < n; i++ {
		c := uint64(t[i])
		h1[i+1] = (h1[i]*base1 + c) % mod1
		p1[i+1] = (p1[i] * base1) % mod1

		h2[i+1] = h2[i]*base2 + c
		p2[i+1] = p2[i] * base2
	}

	// Identify valid start positions for sBegin
	// starts[i] is true if t starts with sBegin at index i
	starts := make([]bool, n+1)
	for i := 0; i <= n-nb; i++ {
		if t[i:i+nb] == sBegin {
			starts[i] = true
		}
	}

	// Identify valid end positions for sEnd
	// ends[j] is true if t ends with sEnd at index j (exclusive)
	ends := make([]bool, n+1)
	for i := ne; i <= n; i++ {
		if t[i-ne:i] == sEnd {
			ends[i] = true
		}
	}

	// Determine minimum length of valid substring
	minLen := nb
	if ne > minLen {
		minLen = ne
	}

	// Store unique hashes
	type pair struct {
		h1, h2 uint64
	}
	seen := make(map[pair]struct{})

	// Iterate over all valid start indices
	for i := 0; i <= n; i++ {
		if !starts[i] {
			continue
		}
		// Iterate over all valid end indices
		for j := i + minLen; j <= n; j++ {
			if !ends[j] {
				continue
			}

			// Compute hash of t[i:j]
			// len = j - i
			
			// Hash 1
			val1 := (h1[j] + mod1 - (h1[i]*p1[j-i])%mod1) % mod1
			
			// Hash 2
			val2 := h2[j] - h2[i]*p2[j-i]

			seen[pair{val1, val2}] = struct{}{}
		}
	}

	fmt.Println(len(seen))
}