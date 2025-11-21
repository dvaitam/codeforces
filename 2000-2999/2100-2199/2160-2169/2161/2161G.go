package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const (
	maxBits = 20
	limit   = 1 << maxBits
	maxC    = maxBits // include synthetic bit 20
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	freq := make([]int, limit)
	seen := make([]bool, limit)
	values := make([]int, 0, n)
	origAND := -1
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		freq[x]++
		if !seen[x] {
			seen[x] = true
			values = append(values, x)
		}
		if i == 0 {
			origAND = x
		} else {
			origAND &= x
		}
	}

	countZero := make([]int, maxC+1)
	for _, v := range values {
		cnt := freq[v]
		for b := 0; b < maxBits; b++ {
			if ((v >> b) & 1) == 0 {
				countZero[b] += cnt
			}
		}
	}
	countZero[maxC] = n

	contain := make([]int, limit)
	copy(contain, freq)
	for bit := 0; bit < maxBits; bit++ {
		step := 1 << uint(bit)
		jump := step << 1
		for start := 0; start < limit; start += jump {
			for offset := 0; offset < step; offset++ {
				idx := start + offset
				idx2 := idx + step
				contain[idx] += contain[idx2]
			}
		}
	}

	cntMasks := make([][]int, maxBits)
	sumMasks := make([][]int64, maxBits)
	for b := 0; b < maxBits; b++ {
		dim := maxBits - b - 1
		size := 1 << uint(dim)
		cntArr := make([]int, size)
		sumArr := make([]int64, size)
		lowMask := (1 << uint(b)) - 1
		shift := b + 1
		for _, v := range values {
			if ((v >> b) & 1) == 0 {
				idx := v >> shift
				cntArr[idx] += freq[v]
				sumArr[idx] += int64(freq[v]) * int64(v&lowMask)
			}
		}
		for bit := 0; bit < dim; bit++ {
			step := 1 << uint(bit)
			jump := step << 1
			for start := 0; start < size; start += jump {
				for offset := 0; offset < step; offset++ {
					idx := start + offset
					idx2 := idx + step
					cntArr[idx] += cntArr[idx2]
					sumArr[idx] += sumArr[idx2]
				}
			}
		}
		cntMasks[b] = cntArr
		sumMasks[b] = sumArr
	}

	bestGood := buildBestGood(values)
	bestBad := buildBestBad(values)

	constInf := int64(1) << 62

	for ; q > 0; q-- {
		var X int
		fmt.Fscan(in, &X)
		hi := highestBit(X)
		total := int64(0)
		for b := 0; b < maxBits; b++ {
			if ((X >> b) & 1) == 0 {
				continue
			}
			dim := maxBits - b - 1
			mask := 0
			if dim > 0 {
				mask = (X >> (b + 1)) & ((1 << uint(dim)) - 1)
			}
			cnt := cntMasks[b][mask]
			sumVal := sumMasks[b][mask]
			term := int64(cnt)*(int64(1<<uint(b))+int64(X&((1<<uint(b))-1))) - sumVal
			total += term
		}

		maskLow := 0
		if hi >= 0 {
			maskLow = (1 << uint(hi+1)) - 1
		}
		highPart := origAND & ^maskLow
		lowPart := 0
		if hi >= 0 {
			if contain[X] == n {
				lowPart = origAND & maskLow
			} else {
				lowPart = X
			}
		}
		baseAND := highPart | lowPart
		if baseAND == X {
			fmt.Fprintln(out, total)
			continue
		}

		bestExtra := constInf
		for b := maxC; b >= hi+1; b-- {
			if countZero[b] == 0 {
				break
			}
			if countZero[b] < 2 {
				continue
			}
			limitMask := (1 << uint(b)) - 1
			constVal := int64(1<<uint(b)) + int64(X&limitMask)
			maxRem := int64(-1)
			if bestGood[b][X] >= 0 {
				maxRem = int64(bestGood[b][X])
			}
			if hi >= 0 && X != 0 {
				arr := bestBad[hi][b]
				if arr != nil {
					tmp := X
					for tmp != 0 {
						lowest := tmp & -tmp
						bitIdx := bits.TrailingZeros(uint(lowest))
						tmp -= lowest
						if bitIdx <= hi {
							if arr[bitIdx] >= 0 {
								rem := int64(arr[bitIdx]) + int64(X)
								if rem > maxRem {
									maxRem = rem
								}
							}
						}
					}
				}
			}
			if maxRem >= 0 {
				extra := constVal - maxRem
				if extra < bestExtra {
					bestExtra = extra
				}
			}
		}
		if bestExtra == constInf {
			bestExtra = 0
		}
		fmt.Fprintln(out, total+bestExtra)
	}
}

func highestBit(x int) int {
	if x == 0 {
		return -1
	}
	return bits.Len(uint(x)) - 1
}

func buildBestGood(values []int) [][]int {
	best := make([][]int, maxC+1)
	for c := 0; c <= maxC; c++ {
		arr := make([]int, limit)
		for i := range arr {
			arr[i] = -1
		}
		mask := (1 << uint(c)) - 1
		for _, v := range values {
			if ((v >> c) & 1) == 0 {
				rem := v & mask
				if rem > arr[v] {
					arr[v] = rem
				}
			}
		}
		for bit := 0; bit < maxBits; bit++ {
			step := 1 << uint(bit)
			jump := step << 1
			for start := 0; start < limit; start += jump {
				for offset := 0; offset < step; offset++ {
					idx := start + offset
					idx2 := idx + step
					if arr[idx2] > arr[idx] {
						arr[idx] = arr[idx2]
					}
				}
			}
		}
		best[c] = arr
	}
	return best
}

func buildBestBad(values []int) [][][]int {
	best := make([][][]int, maxBits)
	midMasks := make([][]int, maxBits)
	for hi := 0; hi < maxBits; hi++ {
		arrays := make([][]int, maxC+1)
		mids := make([]int, maxC+1)
		for b := hi + 1; b <= maxC; b++ {
			mids[b] = ((1 << uint(b)) - 1) & ^((1 << uint(hi+1)) - 1)
			arr := make([]int, hi+1)
			for i := range arr {
				arr[i] = -1
			}
			arrays[b] = arr
		}
		best[hi] = arrays
		midMasks[hi] = mids
	}

	for _, v := range values {
		for b := 0; b <= maxC; b++ {
			if ((v >> b) & 1) != 0 {
				continue
			}
			limitHi := b
			if limitHi > maxBits {
				limitHi = maxBits
			}
			for hi := 0; hi < limitHi; hi++ {
				arr := best[hi][b]
				if arr == nil {
					continue
				}
				midVal := v & midMasks[hi][b]
				for t := 0; t <= hi; t++ {
					if ((v >> t) & 1) == 0 {
						if midVal > arr[t] {
							arr[t] = midVal
						}
					}
				}
			}
		}
	}
	return best
}
