package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	// frequency of strings by length and total sum
	var freq [6][46]int
	// prefix count: length, prefixLen (<length), prefixSum, totalSum
	var pref [6][6][46][46]int

	strs := make([]string, n)
	lens := make([]int, n)
	sums := make([]int, n)
	prefixes := make([][]int, n)

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		strs[i] = s
		l := len(s)
		lens[i] = l
		prefix := make([]int, l+1)
		sum := 0
		for j := 0; j < l; j++ {
			d := int(s[j] - '0')
			sum += d
			prefix[j+1] = sum
			if j+1 < l {
				pref[l][j+1][prefix[j+1]][sum]++
			}
		}
		sums[i] = sum
		prefixes[i] = prefix
		freq[l][sum]++
	}

	ans := 0
	for idx := 0; idx < n; idx++ {
		lS := lens[idx]
		sumS := sums[idx]
		preS := prefixes[idx]
		for lenT := 1; lenT <= 5; lenT++ {
			L := lS + lenT
			if L%2 == 1 {
				continue
			}
			half := L / 2
			if half <= lS {
				d := 2*preS[half] - sumS
				if d >= lenT && d <= 9*lenT && d >= 0 && d < 46 {
					ans += freq[lenT][d]
				}
			} else {
				prefLen := half - lS
				if prefLen >= lenT {
					continue
				}
				maxPrefSum := 9 * prefLen
				if maxPrefSum > 45 {
					maxPrefSum = 45
				}
				for ps := 1; ps <= maxPrefSum; ps++ {
					sumT := sumS + 2*ps
					if sumT < lenT || sumT > 9*lenT || sumT >= 46 {
						continue
					}
					ans += pref[lenT][prefLen][ps][sumT]
				}
			}
		}
	}

	fmt.Fprintln(out, ans)
}
