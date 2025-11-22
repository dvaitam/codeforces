package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxA = 100000

// buildSPF returns the smallest prime factor for every number up to limit.
func buildSPF(limit int) []int {
	spf := make([]int, limit+1)
	for i := 2; i <= limit; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	spf[1] = 1
	return spf
}

// divisors enumerates all divisors of x using its prime factorization.
func divisors(x int, spf []int) []int {
	res := []int{1}
	for x > 1 {
		p := spf[x]
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		curLen := len(res)
		mul := 1
		for i := 1; i <= cnt; i++ {
			mul *= p
			for j := 0; j < curLen; j++ {
				res = append(res, res[j]*mul)
			}
		}
	}
	return res
}

// nextDivPosition finds the smallest position in [pos, r] where a value divides k
// (i.e., is one of the provided divisors). Returns r+1 if none found.
func nextDivPosition(posLists [][]int, divs []int, pos, r int) int {
	best := r + 1
	for _, d := range divs {
		lst := posLists[d]
		if len(lst) == 0 {
			continue
		}
		idx := sort.SearchInts(lst, pos)
		if idx < len(lst) {
			if v := lst[idx]; v < best {
				best = v
			}
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	spf := buildSPF(maxA)

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)

		arr := make([]int, n)
		posLists := make([][]int, maxA+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			posLists[arr[i]] = append(posLists[arr[i]], i+1) // store 1-based positions
		}

		for ; q > 0; q-- {
			var k, l, r int
			fmt.Fscan(in, &k, &l, &r)

			curPos := l
			curK := k
			var ans int64

			for curPos <= r {
				if curK == 1 {
					ans += int64(r - curPos + 1)
					break
				}

				divs := divisors(curK, spf)
				nextPos := nextDivPosition(posLists, divs, curPos, r)

				if nextPos > r { // no more divisors in range
					ans += int64(r-curPos+1) * int64(curK)
					break
				}

				if nextPos > curPos {
					ans += int64(nextPos-curPos) * int64(curK)
				}

				val := arr[nextPos-1]
				for curK%val == 0 {
					curK /= val
				}
				ans += int64(curK)
				curPos = nextPos + 1
			}

			fmt.Fprintln(out, ans)
		}
	}
}
