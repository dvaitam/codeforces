package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxVal = 500000 + 5

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m, k, n, s int
	if _, err := fmt.Fscan(in, &m, &k, &n, &s); err != nil {
		return
	}

	a := make([]int, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &a[i])
	}

	need := make([]int, maxVal)
	unique := []int{}
	totalNeed := 0
	for i := 0; i < s; i++ {
		var b int
		fmt.Fscan(in, &b)
		if need[b] == 0 {
			unique = append(unique, b)
		}
		need[b]++
		totalNeed++
	}

	freq := make([]int, maxVal)
	l := 1
	matched := 0
	remTotal := m - n*k
	found := false
	remNeed := make([]int, maxVal)

	var res []int

	for r := 1; r <= m && !found; r++ {
		val := a[r]
		freq[val]++
		if need[val] > 0 && freq[val] <= need[val] {
			matched++
		}
		for matched == totalNeed && !found {
			lenSeg := r - l + 1
			if lenSeg >= k {
				keepBlocks := min((l-1)/k, n-1)
				keepPrefix := keepBlocks * k
				remPrefix := (l - 1) - keepPrefix
				remInside := lenSeg - k
				if remPrefix+remInside <= remTotal {
					res = res[:0]
					for i := keepPrefix + 1; i <= l-1; i++ {
						res = append(res, i)
					}
					for _, t := range unique {
						remNeed[t] = need[t]
					}
					toRemove := remInside
					for i := l; i <= r; i++ {
						val := a[i]
						if remNeed[val] > 0 {
							remNeed[val]--
							continue
						}
						if toRemove > 0 {
							res = append(res, i)
							toRemove--
						}
					}
					fmt.Fprintln(out, len(res))
					if len(res) > 0 {
						for i, v := range res {
							if i > 0 {
								fmt.Fprint(out, " ")
							}
							fmt.Fprint(out, v)
						}
						fmt.Fprintln(out)
					} else {
						fmt.Fprintln(out)
					}
					found = true
					break
				}
			}
			leftVal := a[l]
			freq[leftVal]--
			if need[leftVal] > 0 && freq[leftVal] < need[leftVal] {
				matched--
			}
			l++
		}
	}

	if !found {
		fmt.Fprintln(out, -1)
	}
}
