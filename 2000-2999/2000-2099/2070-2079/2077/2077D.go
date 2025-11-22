package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	c, err := fs.r.ReadByte()
	for err == nil && (c == ' ' || c == '\n' || c == '\r' || c == '\t') {
		c, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for err == nil && c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = fs.nextInt()
		}

		// Find largest possible maximum side length M such that
		// all elements <= M give total sum > 2*M and count >= 3.
		sorted := append([]int(nil), a...)
		sort.Ints(sorted)
		cand := -1
		var count int
		var sum int64
		for i := 0; i < n; {
			j := i
			for j < n && sorted[j] == sorted[i] {
				j++
			}
			val := sorted[i]
			cnt := j - i
			count += cnt
			sum += int64(val * cnt)
			if count >= 3 && sum > 2*int64(val) {
				cand = val
			}
			i = j
		}

		if cand == -1 {
			fmt.Fprintln(out, -1)
			continue
		}

		M := cand
		// Filter elements that do not exceed M.
		filtered := make([]int, 0, n)
		for _, v := range a {
			if v <= M {
				filtered = append(filtered, v)
			}
		}

		m := len(filtered)
		sufSum := make([]int64, m+1)
		for i := m - 1; i >= 0; i-- {
			sufSum[i] = sufSum[i+1] + int64(filtered[i])
		}

		target := 2*int64(M) + 1 // need strictly greater than 2*M
		pos := 0
		curSum := int64(0)
		curLen := 0
		res := make([]int, 0, m)

		// deque for sliding window maximum within current feasible range
		deque := make([]int, 0)
		rAdded := pos - 1

		for pos < m {
			if curLen >= 3 && curSum > 2*int64(M) {
				// Condition already satisfied: appending remaining elements makes sequence lexicographically larger.
				res = append(res, filtered[pos:]...)
				break
			}

			// Determine farthest feasible index where we can still achieve the goal.
			lo, hi := pos, m-1
			ans := pos - 1
			for lo <= hi {
				mid := (lo + hi) >> 1
				if curSum+sufSum[mid] >= target && curLen+(m-mid) >= 3 {
					ans = mid
					lo = mid + 1
				} else {
					hi = mid - 1
				}
			}

			if ans < pos {
				// Should not happen since M was feasible.
				res = nil
				break
			}

			// Extend deque window to [pos, ans]
			for i := rAdded + 1; i <= ans; i++ {
				val := filtered[i]
				for len(deque) > 0 && filtered[deque[len(deque)-1]] < val {
					deque = deque[:len(deque)-1]
				}
				deque = append(deque, i)
			}
			rAdded = ans
			for len(deque) > 0 && deque[0] < pos {
				deque = deque[1:]
			}

			idx := deque[0]
			res = append(res, filtered[idx])
			curSum += int64(filtered[idx])
			curLen++
			pos = idx + 1
			for len(deque) > 0 && deque[0] < pos {
				deque = deque[1:]
			}
		}

		if res == nil {
			fmt.Fprintln(out, -1)
			continue
		}

		fmt.Fprintln(out, len(res))
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
