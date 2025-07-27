package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func smallestForLen(length int, digits []int) string {
	sort.Ints(digits)
	// choose first non-zero digit for the leading position
	first := -1
	for _, d := range digits {
		if d != 0 {
			first = d
			break
		}
	}
	if first == -1 {
		// cannot form a positive number
		return ""
	}
	res := make([]byte, length)
	res[0] = byte('0' + first)
	fill := digits[0]
	for i := 1; i < length; i++ {
		res[i] = byte('0' + fill)
	}
	return string(res)
}

func attemptSameLen(n string, digits []int) (string, bool) {
	sort.Ints(digits)
	L := len(n)
	res := make([]byte, L)
	var dfs func(pos int, tight bool) bool
	dfs = func(pos int, tight bool) bool {
		if pos == L {
			return true
		}
		nd := int(n[pos] - '0')
		for _, d := range digits {
			if pos == 0 && d == 0 {
				continue
			}
			if tight {
				if d < nd {
					continue
				}
				res[pos] = byte('0' + d)
				if d == nd {
					if dfs(pos+1, true) {
						return true
					}
				} else { // d > nd
					for i := pos + 1; i < L; i++ {
						res[i] = byte('0' + digits[0])
					}
					return true
				}
			} else {
				res[pos] = byte('0' + d)
				for i := pos + 1; i < L; i++ {
					res[i] = byte('0' + digits[0])
				}
				return true
			}
		}
		return false
	}
	if dfs(0, true) {
		return string(res), true
	}
	return "", false
}

func buildCandidate(n string, digits []int) string {
	if cand, ok := attemptSameLen(n, digits); ok {
		return cand
	}
	return smallestForLen(len(n)+1, digits)
}

func cmpLess(a, b string) bool {
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		s := strconv.Itoa(n)
		best := ""
		if k == 1 {
			for d := 1; d <= 9; d++ {
				cand := buildCandidate(s, []int{d})
				if best == "" || cmpLess(cand, best) {
					best = cand
				}
			}
		} else { // k == 2
			for i := 0; i <= 9; i++ {
				for j := i; j <= 9; j++ {
					digits := []int{i}
					if j != i {
						digits = append(digits, j)
					}
					if len(digits) == 1 && digits[0] == 0 {
						continue
					}
					cand := buildCandidate(s, digits)
					if cand == "" {
						continue
					}
					if best == "" || cmpLess(cand, best) {
						best = cand
					}
				}
			}
		}
		fmt.Fprintln(writer, best)
	}
}
