package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type solver struct {
	allVals       []int64
	allowedVals   []int64
	prefixAllowed []int64
	required      int
}

func lowerBound(arr []int64, target int64) int {
	return sort.Search(len(arr), func(i int) bool { return arr[i] >= target })
}

func (s *solver) costForMedian(x int64) (int64, int, bool) {
	total := len(s.allVals)
	idx := lowerBound(s.allVals, x)
	cntGE := total - idx
	need := s.required - cntGE
	if need <= 0 {
		return 0, 0, true
	}
	pos := lowerBound(s.allowedVals, x)
	if need > pos {
		return 0, need, false
	}
	sumVals := s.prefixAllowed[pos] - s.prefixAllowed[pos-need]
	cost := int64(need)*x - sumVals
	return cost, need, true
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	const negInf = int64(-1 << 60)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		allVals := append([]int64(nil), a...)
		sort.Slice(allVals, func(i, j int) bool { return allVals[i] < allVals[j] })

		allowedVals := make([]int64, 0)
		prefix := []int64{0}
		hasAllowed := false
		maxAllowed := negInf
		maxDisallowed := negInf
		for i := 0; i < n; i++ {
			if b[i] == 1 {
				hasAllowed = true
				if a[i] > maxAllowed {
					maxAllowed = a[i]
				}
				allowedVals = append(allowedVals, a[i])
			} else {
				if a[i] > maxDisallowed {
					maxDisallowed = a[i]
				}
			}
		}
		if !hasAllowed {
			maxAllowed = negInf
		}
		sort.Slice(allowedVals, func(i, j int) bool { return allowedVals[i] < allowedVals[j] })
		for _, v := range allowedVals {
			prefix = append(prefix, prefix[len(prefix)-1]+v)
		}

		m := n / 2
		required := n - m + 1
		s := solver{
			allVals:       allVals,
			allowedVals:   allowedVals,
			prefixAllowed: prefix,
			required:      required,
		}

		candSet := make(map[int64]struct{})
		for _, val := range a {
			if val >= 1 {
				candSet[val] = struct{}{}
			}
			if val+1 >= 1 {
				candSet[val+1] = struct{}{}
			}
		}
		maxAll := allVals[len(allVals)-1]
		limitHigh := maxAll
		if hasAllowed {
			if maxAllowed+k > limitHigh {
				limitHigh = maxAllowed + k
			}
		}
		candSet[limitHigh+1] = struct{}{}

		points := make([]int64, 0, len(candSet))
		for v := range candSet {
			points = append(points, v)
		}
		sort.Slice(points, func(i, j int) bool { return points[i] < points[j] })

		best := int64(0)
		eval := func(x int64, cost int64, needPositive bool) {
			if x < 1 || cost > k {
				return
			}
			kRem := k - cost
			maxVal := maxDisallowed
			if hasAllowed {
				baseAllowed := maxAllowed
				if needPositive && x > baseAllowed {
					baseAllowed = x
				}
				allowedVal := baseAllowed + kRem
				if allowedVal > maxVal {
					maxVal = allowedVal
				}
			}
			if maxVal < negInf/2 {
				return
			}
			score := x + maxVal
			if score > best {
				best = score
			}
		}

		for i := 0; i+1 < len(points); i++ {
			L := points[i]
			R := points[i+1]
			if L >= R {
				continue
			}
			costL, need, ok := s.costForMedian(L)
			if !ok || costL > k {
				continue
			}
			eval(L, costL, need > 0)
			if need == 0 {
				Xcand := R - 1
				if Xcand > L {
					eval(Xcand, costL, false)
				}
				continue
			}
			maxDelta := (k - costL) / int64(need)
			if maxDelta <= 0 {
				continue
			}
			Xcand := L + maxDelta
			if Xcand >= R {
				Xcand = R - 1
			}
			if Xcand > L {
				costX := costL + int64(need)*(Xcand-L)
				eval(Xcand, costX, true)
			}
		}

		fmt.Fprintln(out, best)
	}
}
