package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func attempt(n int, segs [][]int, start int, first, second int) ([]int, bool) {
	used := make([]bool, len(segs))
	used[start] = true
	prefixSet := make(map[int]bool)
	prefixSet[first] = true
	prefixSet[second] = true
	perm := []int{first, second}
	for len(perm) < n {
		foundIdx := -1
		newVal := 0
		for i, s := range segs {
			if used[i] {
				continue
			}
			diffCount := 0
			candidate := 0
			for _, v := range s {
				if !prefixSet[v] {
					diffCount++
					candidate = v
					if diffCount > 1 {
						break
					}
				}
			}
			if diffCount == 1 {
				foundIdx = i
				newVal = candidate
				break
			}
		}
		if foundIdx == -1 {
			return nil, false
		}
		used[foundIdx] = true
		prefixSet[newVal] = true
		perm = append(perm, newVal)
	}
	for _, u := range used {
		if !u {
			return nil, false
		}
	}
	return perm, true
}

func verify(perm []int, segs [][]int) bool {
	n := len(perm)
	// map serialized segment -> count
	counter := make(map[string]int)
	for _, s := range segs {
		key := fmt.Sprint(s)
		counter[key]++
	}
	for r := 2; r <= n; r++ {
		found := false
		for l := 1; l < r; l++ {
			sub := append([]int(nil), perm[l-1:r]...)
			sort.Ints(sub)
			key := fmt.Sprint(sub)
			if counter[key] > 0 {
				counter[key]--
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func solve(n int, segs [][]int) []int {
	for idx, s := range segs {
		if len(s) == 2 {
			a, b := s[0], s[1]
			if perm, ok := attempt(n, segs, idx, a, b); ok {
				if verify(perm, segs) {
					return perm
				}
			}
			if perm, ok := attempt(n, segs, idx, b, a); ok {
				if verify(perm, segs) {
					return perm
				}
			}
		}
	}
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		segs := make([][]int, n-1)
		for i := 0; i < n-1; i++ {
			var k int
			fmt.Fscan(reader, &k)
			seg := make([]int, k)
			for j := 0; j < k; j++ {
				fmt.Fscan(reader, &seg[j])
			}
			segs[i] = seg
		}
		perm := solve(n, segs)
		for i, v := range perm {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
