package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxA = 100000

var spf [maxA + 1]int

func initSieve() {
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxA; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func factorize(x int) []int {
	var res []int
	for x > 1 {
		p := spf[x]
		res = append(res, p)
		for x%p == 0 {
			x /= p
		}
	}
	return res
}

func contains(list []int, x int) bool {
	for _, v := range list {
		if v == x {
			return true
		}
	}
	return false
}

func containsAll(list []int, must []int) bool {
	for _, v := range must {
		if !contains(list, v) {
			return false
		}
	}
	return true
}

func intersect(a, b []int) []int {
	if len(a) > len(b) {
		a, b = b, a
	}
	var res []int
	for _, v := range a {
		if contains(b, v) {
			res = append(res, v)
		}
	}
	return res
}

type dfsState struct {
	k          int
	factors    [][]int
	isHeavy    []bool
	pos        map[int][]int
	used       []bool
	missingSet map[int]bool
	goodSet    map[int]bool
	ans        []int
}

func (st *dfsState) dfs(selected []int, missing []int, common []int) bool {
	if len(selected) == st.k {
		st.ans = append([]int{}, selected...)
		return true
	}
	if len(selected)+len(common) < st.k {
		return false
	}

	// pick anchor prime with smallest candidate list
	anchor := missing[0]
	minSize := len(st.pos[anchor])
	for _, p := range missing {
		if len(st.pos[p]) < minSize {
			minSize = len(st.pos[p])
			anchor = p
		}
	}

	for _, idx := range st.pos[anchor] {
		if st.used[idx] || !st.isHeavy[idx] {
			continue
		}
		if !containsAll(st.factors[idx], missing) {
			continue
		}
		var candidates []int
		for _, p := range common {
			if st.missingSet[p] {
				continue
			}
			if !contains(st.factors[idx], p) {
				candidates = append(candidates, p)
			}
		}
		for _, miss := range candidates {
			newCommon := intersect(common, st.factors[idx])
			if len(selected)+1+len(newCommon) < st.k {
				continue
			}
			st.used[idx] = true
			st.missingSet[miss] = true
			if st.dfs(append(selected, idx), append(missing, miss), newCommon) {
				return true
			}
			st.missingSet[miss] = false
			st.used[idx] = false
		}
	}
	return false
}

func findMaxIndependent(a []int) []int {
	n := len(a)
	factors := make([][]int, n)
	freq := make([]int, maxA+1)
	for i, v := range a {
		factors[i] = factorize(v)
		for _, p := range factors[i] {
			freq[p]++
		}
	}

	// positions for all primes
	posAll := make(map[int][]int)
	for idx, f := range factors {
		for _, p := range f {
			posAll[p] = append(posAll[p], idx)
		}
	}

	maxK := 7
	if n < maxK {
		maxK = n
	}

	// handle k=2 and k=1 fallback later
	for k := maxK; k >= 3; k-- {
		good := make([]bool, maxA+1)
		for p, c := range freq {
			if c >= k-1 {
				good[p] = true
			}
		}

		// build good factors and heavy numbers
		goodFactors := make([][]int, n)
		isHeavy := make([]bool, n)
		heavyCount := 0
		for i := 0; i < n; i++ {
			for _, p := range factors[i] {
				if good[p] {
					goodFactors[i] = append(goodFactors[i], p)
				}
			}
			if len(goodFactors[i]) >= k-1 {
				isHeavy[i] = true
				heavyCount++
			}
		}
		if heavyCount < k {
			continue
		}

		// heavy frequency and positions restricted to heavy numbers
		heavyFreq := make([]int, maxA+1)
		pos := make(map[int][]int)
		for i := 0; i < n; i++ {
			if !isHeavy[i] {
				continue
			}
			for _, p := range factors[i] {
				heavyFreq[p]++
				pos[p] = append(pos[p], i)
			}
		}

		// list of primes that appear enough among heavy numbers
		var heavyPrimes []int
		for p, c := range heavyFreq {
			if c >= k-1 {
				heavyPrimes = append(heavyPrimes, p)
			}
		}
		if len(heavyPrimes) < k {
			continue
		}

		used := make([]bool, n)
		missingSet := make(map[int]bool)
		state := dfsState{
			k:          k,
			factors:    goodFactors,
			isHeavy:    isHeavy,
			pos:        pos,
			used:       used,
			missingSet: missingSet,
		}

		for i := 0; i < n; i++ {
			if !isHeavy[i] || len(goodFactors[i]) < k-1 {
				continue
			}
			common := goodFactors[i]
			if len(common)+1 < k {
				continue
			}

			for _, p := range heavyPrimes {
				if heavyFreq[p] < k-1 {
					continue
				}
				if contains(factors[i], p) {
					continue
				}
				// start DFS
				for key := range state.missingSet {
					delete(state.missingSet, key)
				}
				for idx := range state.used {
					state.used[idx] = false
				}
				state.used[i] = true
				state.missingSet[p] = true
				selected := []int{i}
				missing := []int{p}
				state.factors = goodFactors
				if state.dfs(selected, missing, common) {
					return state.ans
				}
				state.used[i] = false
				state.missingSet[p] = false
			}
		}
	}

	// try to find independent pair
	if n >= 2 {
		type pair struct {
			val int
			idx int
		}
		arr := make([]pair, n)
		for i, v := range a {
			arr[i] = pair{v, i}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
		for i := 0; i+1 < n; i++ {
			if arr[i+1].val%arr[i].val != 0 {
				return []int{arr[i].idx, arr[i+1].idx}
			}
		}
	}

	return []int{0}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	initSieve()

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		indices := findMaxIndependent(a)
		fmt.Fprintln(out, len(indices))
		for i, idx := range indices {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[idx])
		}
		fmt.Fprintln(out)
	}
}
