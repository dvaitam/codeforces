package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const limit = 1_000_000_000

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		a := make([]int, n)
		b := make([]int, n)
		maxA, maxB := 0, 0
		var sumA, sumB int64

		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxA {
				maxA = a[i]
			}
			sumA += int64(a[i])
		}

		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
			if b[i] > maxB {
				maxB = b[i]
			}
			sumB += int64(b[i])
		}

		if maxB > maxA || sumA < sumB {
			fmt.Fprintln(out, -1)
			continue
		}

		countsB := make(map[int]int, n)
		for _, val := range b {
			countsB[val]++
		}

		diff := sumA - sumB
		if diff == 0 {
			if sameMultiset(a, b) {
				ans := maxA + 1
				if ans > limit {
					ans = limit
				}
				fmt.Fprintln(out, ans)
			} else {
				fmt.Fprintln(out, -1)
			}
			continue
		}

		candidates := divisorsFiltered(diff, maxB)
		answer := -1
		for _, k64 := range candidates {
			if k64 > limit {
				continue
			}
			k := int(k64)
			if checkCandidate(a, countsB, k) {
				answer = k
				break
			}
		}

		fmt.Fprintln(out, answer)
	}
}

func sameMultiset(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	copyA := append([]int(nil), a...)
	copyB := append([]int(nil), b...)
	sort.Ints(copyA)
	sort.Ints(copyB)
	for i := range copyA {
		if copyA[i] != copyB[i] {
			return false
		}
	}
	return true
}

func divisorsFiltered(v int64, maxB int) []int64 {
	res := make([]int64, 0)
	maxB64 := int64(maxB)
	for i := int64(1); i*i <= v; i++ {
		if v%i != 0 {
			continue
		}
		d1 := i
		d2 := v / i
		if d1 > maxB64 && d1 <= limit {
			res = append(res, d1)
		}
		if d2 != d1 && d2 > maxB64 && d2 <= limit {
			res = append(res, d2)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res
}

func checkCandidate(a []int, countsB map[int]int, k int) bool {
	temp := make(map[int]int, len(countsB))
	for _, val := range a {
		r := val % k
		need, ok := countsB[r]
		if !ok {
			return false
		}
		temp[r]++
		if temp[r] > need {
			return false
		}
	}
	for key, val := range countsB {
		if temp[key] != val {
			return false
		}
	}
	return true
}
