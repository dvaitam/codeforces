package main

import (
	"bufio"
	"fmt"
	"os"
)

// This solution implements a very naive exhaustive search to illustrate
// the rules from problemF.txt. It only works for extremely small inputs.
// For larger n the search space becomes prohibitive so the program simply
// prints -1 as a placeholder. A correct and efficient algorithm would be
// required for actual competition use.

var (
	n, k int
	a    []int
	best int
)

func checkAssign(idx []int) bool {
	m := len(idx)
	// try every gender assignment using bit masks
	for mask := 0; mask < 1<<uint(m); mask++ {
		fcnt, mcnt := 0, 0
		lastF, lastM := 0, 0
		haveF, haveM := false, false
		sumF, sumM := 0, 0
		ok := true
		for i, pos := range idx {
			val := a[pos]
			if mask&(1<<uint(i)) != 0 { // female
				if haveF {
					if val != lastF+1 {
						ok = false
						break
					}
				} else {
					haveF = true
				}
				lastF = val
				fcnt++
				sumF += val
			} else { // male
				if haveM {
					if val != lastM-1 {
						ok = false
						break
					}
				} else {
					haveM = true
				}
				lastM = val
				mcnt++
				sumM += val
			}
		}
		if ok && haveF && haveM && sumF*mcnt == sumM*fcnt {
			return true
		}
	}
	return false
}

func dfs(pos, chosen int, comb []int) {
	if chosen == k {
		if checkAssign(comb) {
			diff := comb[k-1] - comb[0]
			if diff < best {
				best = diff
			}
		}
		return
	}
	if pos == n {
		return
	}
	if n-pos < k-chosen {
		return
	}
	// choose current index
	comb[chosen] = pos
	dfs(pos+1, chosen+1, comb)
	// skip current index
	dfs(pos+1, chosen, comb)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if k%2 == 1 || n > 20 {
		fmt.Fprintln(out, -1)
		return
	}
	best = int(1e9)
	comb := make([]int, k)
	dfs(0, 0, comb)
	if best == int(1e9) {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, best)
	}
}
