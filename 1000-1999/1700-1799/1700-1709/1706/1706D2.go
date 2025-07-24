package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func solveCase(n, k int, a []int) int {
	sort.Ints(a)
	minA := a[0]
	maxA := a[n-1]
	flags := make([]bool, minA+1)
	flags[0] = true
	for _, x := range a {
		p := 1
		for p <= k {
			val := x / p
			if val <= minA {
				flags[val] = true
			}
			if val == 0 {
				break
			}
			p = x/val + 1
		}
	}
	// gather candidates in descending order
	cands := make([]int, 0)
	for v := minA; v >= 0; v-- {
		if flags[v] {
			cands = append(cands, v)
		}
	}
	if len(cands) == 0 {
		cands = append(cands, 0)
	}

	pi := make([]int, n)
	hi := make([]int, n)
	bucket := make([][]int, maxA+1)
	freq := make([]int, maxA+1)

	start := cands[0]
	for i, x := range a {
		if start == 0 {
			pi[i] = k
		} else {
			pi[i] = x / start
			if pi[i] > k {
				pi[i] = k
			}
			if pi[i] < 1 {
				pi[i] = 1
			}
		}
		hi[i] = x / pi[i]
		freq[hi[i]]++
		t := x / (pi[i] + 1)
		bucket[t] = append(bucket[t], i)
	}

	hiMax := maxA
	for hiMax > 0 && freq[hiMax] == 0 {
		hiMax--
	}
	best := hiMax - start
	curr := start

	for idx := 1; idx < len(cands); idx++ {
		l := cands[idx]
		for t := curr - 1; t >= l && t >= 0; t-- {
			b := bucket[t]
			if len(b) > 0 {
				bucket[t] = nil
				for _, id := range b {
					var newPi int
					if l == 0 {
						newPi = k
					} else {
						newPi = a[id] / l
						if newPi > k {
							newPi = k
						}
					}
					if newPi > pi[id] {
						freq[hi[id]]--
						pi[id] = newPi
						hi[id] = a[id] / pi[id]
						freq[hi[id]]++
						nt := a[id] / (pi[id] + 1)
						bucket[nt] = append(bucket[nt], id)
					} else {
						// Should not happen
						nt := a[id] / (pi[id] + 1)
						bucket[nt] = append(bucket[nt], id)
					}
				}
			}
		}
		curr = l
		for hiMax > 0 && freq[hiMax] == 0 {
			hiMax--
		}
		diff := hiMax - l
		if diff < best {
			best = diff
		}
	}

	return best
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		res := solveCase(n, k, a)
		fmt.Fprintln(writer, res)
	}
}
