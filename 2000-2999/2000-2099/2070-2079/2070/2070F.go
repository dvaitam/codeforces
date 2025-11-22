package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// The solution treats friend preference masks as polynomials:
//  - for pizzas with odd number of slices, variable x_i is nilpotent (x_i^2 = 0)
//  - for pizzas with even number of slices, variable y_i is idempotent (y_i^2 = y_i)
// A friend liking mask M contributes monomial product of corresponding variables.
// Valid pairs correspond to the square of this polynomial; nilpotent variables
// ensure overlapping odd pizzas vanish (quarrel), idempotent variables model OR on even pizzas.
// We separate odd/even bits. OR convolution on even bits is turned into pointwise
// multiplication via zeta transform. Disjoint-union convolution on odd bits is handled
// with the classic subset convolution in O(k * 2^k), where k = #odd pizzas.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	fullSize := 1 << n

	// Read friend masks and count occurrences.
	countFull := make([]int64, fullSize)
	masks := make([]int, m)
	for i := 0; i < m; i++ {
		var s string
		fmt.Fscan(in, &s)
		mask := 0
		for j := 0; j < len(s); j++ {
			mask |= 1 << (int(s[j]) - 'A')
		}
		masks[i] = mask
		countFull[mask]++
	}

	// Read pizza sizes and determine odd/even sets.
	a := make([]int64, n)
	totalSlices := int64(0)
	oddMask := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		totalSlices += a[i]
		if a[i]%2 == 1 {
			oddMask |= 1 << i
		}
	}
	evenMask := (^oddMask) & (fullSize - 1)

	// Build mapping from original bits to compressed odd/even indices.
	oddIdx := make([]int, n)
	evenIdx := make([]int, n)
	for i := 0; i < n; i++ {
		oddIdx[i] = -1
		evenIdx[i] = -1
	}
	oddBits := make([]int, 0)
	evenBits := make([]int, 0)
	for i := 0; i < n; i++ {
		if oddMask&(1<<i) != 0 {
			oddIdx[i] = len(oddBits)
			oddBits = append(oddBits, i)
		} else {
			evenIdx[i] = len(evenBits)
			evenBits = append(evenBits, i)
		}
	}
	oddCnt := len(oddBits)
	evenCnt := len(evenBits)
	sizeO := 1 << oddCnt
	sizeE := 1 << evenCnt
	size := sizeO * sizeE // equals 1<<n

	// Precompute compressed masks to full masks for reconstruction.
	oddFull := make([]int, sizeO)
	for o := 0; o < sizeO; o++ {
		full := 0
		for i, bit := range oddBits {
			if o&(1<<i) != 0 {
				full |= 1 << bit
			}
		}
		oddFull[o] = full
	}
	evenFull := make([]int, sizeE)
	for e := 0; e < sizeE; e++ {
		full := 0
		for i, bit := range evenBits {
			if e&(1<<i) != 0 {
				full |= 1 << bit
			}
		}
		evenFull[e] = full
	}

	// Base array over compressed odd/even masks.
	base := make([]int64, size)
	selfEven := make([]int64, sizeE) // counts with odd part == 0, used to subtract self-pairs
	for mask, c := range countFull {
		if c == 0 {
			continue
		}
		oPart := 0
		ePart := 0
		tmp := mask
		for bit := 0; bit < n; bit++ {
			if tmp&(1<<bit) == 0 {
				continue
			}
			if oddIdx[bit] != -1 {
				oPart |= 1 << oddIdx[bit]
			} else {
				ePart |= 1 << evenIdx[bit]
			}
		}
		base[oPart*sizeE+ePart] += c
		if oPart == 0 {
			selfEven[ePart] += c
		}
	}

	// OR zeta transform over even dimension for each odd mask.
	if evenCnt > 0 {
		for o := 0; o < sizeO; o++ {
			offset := o * sizeE
			for b := 0; b < evenCnt; b++ {
				step := 1 << b
				for mask := 0; mask < sizeE; mask++ {
					if mask&step != 0 {
						base[offset+mask] += base[offset+mask^step]
					}
				}
			}
		}
	}

	// popcount for odd masks.
	popOdd := make([]int, sizeO)
	for o := 0; o < sizeO; o++ {
		popOdd[o] = bits.OnesCount(uint(o))
	}

	// Layers by popcount for subset convolution on odd bits.
	layers := make([][]int64, oddCnt+1)
	for i := 0; i <= oddCnt; i++ {
		layers[i] = make([]int64, size)
	}
	for o := 0; o < sizeO; o++ {
		pop := popOdd[o]
		copy(layers[pop][o*sizeE:(o+1)*sizeE], base[o*sizeE:(o+1)*sizeE])
	}

	// Zeta transform over odd bits on each layer (subset sums).
	if oddCnt > 0 {
		for p := 0; p <= oddCnt; p++ {
			layer := layers[p]
			for b := 0; b < oddCnt; b++ {
				step := 1 << b
				for o := 0; o < sizeO; o++ {
					if o&step != 0 {
						from := (o ^ step) * sizeE
						to := o * sizeE
						for e := 0; e < sizeE; e++ {
							layer[to+e] += layer[from+e]
						}
					}
				}
			}
		}
	}

	// Multiply in transformed space: for each popcount.
	resLayers := make([][]int64, oddCnt+1)
	for i := 0; i <= oddCnt; i++ {
		resLayers[i] = make([]int64, size)
	}
	for p := 0; p <= oddCnt; p++ {
		res := resLayers[p]
		for idx := 0; idx < size; idx++ {
			var sum int64
			for i := 0; i <= p; i++ {
				sum += layers[i][idx] * layers[p-i][idx]
			}
			res[idx] = sum
		}
	}

	// Inverse zeta over odd bits.
	if oddCnt > 0 {
		for p := 0; p <= oddCnt; p++ {
			layer := resLayers[p]
			for b := 0; b < oddCnt; b++ {
				step := 1 << b
				for o := 0; o < sizeO; o++ {
					if o&step != 0 {
						from := (o ^ step) * sizeE
						to := o * sizeE
						for e := 0; e < sizeE; e++ {
							layer[to+e] -= layer[from+e]
						}
					}
				}
			}
		}
	}

	// Precompute slice sums for all full masks.
	sliceSum := make([]int64, fullSize)
	for mask := 1; mask < fullSize; mask++ {
		lsb := mask & -mask
		bit := bits.TrailingZeros(uint(lsb))
		sliceSum[mask] = sliceSum[mask^lsb] + a[bit]
	}

	// Inverse OR zeta over even bits per odd mask, subtract self pairs, divide by 2, and accumulate answers.
	ans := make([]int64, totalSlices+1)
	for o := 0; o < sizeO; o++ {
		pop := popOdd[o]
		arr := resLayers[pop][o*sizeE : (o+1)*sizeE]
		// Inverse OR transform.
		if evenCnt > 0 {
			for b := 0; b < evenCnt; b++ {
				step := 1 << b
				for e := 0; e < sizeE; e++ {
					if e&step != 0 {
						arr[e] -= arr[e^step]
					}
				}
			}
		}
		for e := 0; e < sizeE; e++ {
			ordered := arr[e]
			if o == 0 {
				ordered -= selfEven[e] // remove (friend, friend) pairs
			}
			if ordered == 0 {
				continue
			}
			unordered := ordered / 2
			if unordered == 0 {
				continue
			}
			fullMask := oddFull[o] | evenFull[e]
			k := totalSlices - sliceSum[fullMask]
			ans[k] += unordered
		}
	}

	// Output.
	fmt.Fprint(out, ans[0])
	for i := 1; i < len(ans); i++ {
		fmt.Fprint(out, " ", ans[i])
	}
	fmt.Fprintln(out)
}
