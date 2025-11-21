package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	k     int64
	n, l  int
	radii []int64
)

type checkInfo struct {
	assignPossible bool
	crossSatisfied bool
	topCount       int
	topSum         int64
	topMax         int64
	topSecond      int64
	constant       int64
}

func lowerBound(arr []int64, target int64) int {
	lo, hi := 0, len(arr)
	for lo < hi {
		mid := (lo + hi) >> 1
		if arr[mid] >= target {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func insertSorted(arr []int64, val int64) []int64 {
	idx := lowerBound(arr, val)
	arr = append(arr, 0)
	copy(arr[idx+1:], arr[idx:])
	arr[idx] = val
	return arr
}

func (info *checkInfo) addTop(val int64) {
	info.topSum += val
	switch info.topCount {
	case 0:
		info.topMax = val
		info.topSecond = val
	case 1:
		if val >= info.topMax {
			info.topSecond = info.topMax
			info.topMax = val
		} else {
			info.topSecond = val
		}
	default:
		if val >= info.topMax {
			info.topSecond = info.topMax
			info.topMax = val
		} else if val > info.topSecond {
			info.topSecond = val
		}
	}
	info.topCount++
}

func checkInt(q int64) checkInfo {
	info := checkInfo{
		assignPossible: true,
		crossSatisfied: true,
	}
	chains := make([]int64, 0, l)
	for i := 0; i < n; i++ {
		needed := radii[i] + q
		idx := lowerBound(chains, needed)
		if idx < len(chains) {
			chains[idx] = radii[i]
			for idx > 0 && chains[idx] < chains[idx-1] {
				chains[idx], chains[idx-1] = chains[idx-1], chains[idx]
				idx--
			}
		} else {
			if i >= l {
				info.assignPossible = false
				return info
			}
			chains = insertSorted(chains, radii[i])
			info.addTop(radii[i])
		}
	}
	if info.topCount == 0 {
		info.assignPossible = false
		return info
	}
	if info.topCount == 1 {
		info.topSecond = info.topMax
	}
	info.constant = 2*info.topSum - (info.topMax + info.topSecond)
	if info.topCount >= 2 {
		rhs := (k - info.constant)
		lhs := int64(info.topCount-1) * q
		if rhs < 0 || lhs > rhs {
			info.crossSatisfied = false
		}
	}
	return info
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &k, &n, &l); err != nil {
		return
	}
	radii = make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &radii[i])
	}

	lo, hi := int64(0), int64(1000000001)
	best := int64(-1)
	for lo <= hi {
		mid := (lo + hi) >> 1
		info := checkInt(mid)
		if info.assignPossible && info.crossSatisfied {
			best = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	if best < 0 {
		fmt.Println(0)
		return
	}

	nextInfo := checkInt(best + 1)
	if !nextInfo.assignPossible {
		fmt.Println(best)
		return
	}
	if nextInfo.crossSatisfied {
		// Should not happen, but guard to avoid infinite loop
		// Move forward to ensure best is maximal
		for nextInfo.assignPossible && nextInfo.crossSatisfied {
			best++
			nextInfo = checkInt(best + 1)
		}
		if !nextInfo.assignPossible {
			fmt.Println(best)
			return
		}
	}
	numerator := k - nextInfo.constant
	denom := int64(nextInfo.topCount - 1)
	if numerator <= 0 {
		fmt.Println(0)
		return
	}
	g := gcd(numerator, denom)
	numerator /= g
	denom /= g
	if denom == 1 {
		fmt.Println(numerator)
	} else {
		fmt.Printf("%d/%d\n", numerator, denom)
	}
}
