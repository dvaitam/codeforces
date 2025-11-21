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

	var mask string
	if _, err := fmt.Fscan(in, &mask); err != nil {
		return
	}
	n := len(mask)
	segments, pre, _ := parseMask(mask)
	g := len(segments)

	if g == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	sumSeg := 0
	sumPre := 0
	for _, v := range segments {
		sumSeg += v
	}
	for _, v := range pre {
		sumPre += v
	}

	tailCandidate := n - sumSeg - sumPre
	if tailCandidate < 0 {
		fmt.Fprintln(out, -1)
		return
	}

	T, ok := chooseT(pre, tailCandidate)
	if !ok {
		fmt.Fprintln(out, -1)
		return
	}

	beforeFirst := pre[0] - T
	if beforeFirst < 0 {
		fmt.Fprintln(out, -1)
		return
	}

	between := make([]int, g-1)
	for i := 1; i < g; i++ {
		between[i-1] = pre[i] - T - 1
		if between[i-1] < 0 {
			fmt.Fprintln(out, -1)
			return
		}
	}

	tailAmount := tailCandidate - T
	if tailAmount < 0 {
		fmt.Fprintln(out, -1)
		return
	}

	var profile []int
	if !addZeroBlocks(&profile, beforeFirst, T) {
		fmt.Fprintln(out, -1)
		return
	}

	for i := 0; i < g; i++ {
		profile = append(profile, segments[i]+T)
		if i+1 < g {
			if !addZeroBlocks(&profile, between[i], T) {
				fmt.Fprintln(out, -1)
				return
			}
		}
	}

	if !addZeroBlocks(&profile, tailAmount, T) {
		fmt.Fprintln(out, -1)
		return
	}

	genMask, ok := buildMask(n, profile)
	if !ok || genMask != mask {
		fmt.Fprintln(out, -1)
		return
	}

	fmt.Fprintln(out, len(profile))
	if len(profile) > 0 {
		for i, v := range profile {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

func parseMask(mask string) ([]int, []int, int) {
	n := len(mask)
	var segments []int
	var pre []int

	i := 0
	gap := 0
	for i < n && mask[i] == '_' {
		gap++
		i++
	}
	if i == n {
		return segments, pre, gap
	}

	for i < n {
		pre = append(pre, gap)
		cnt := 0
		for i < n && mask[i] == '#' {
			cnt++
			i++
		}
		segments = append(segments, cnt)
		gap = 0
		for i < n && mask[i] == '_' {
			gap++
			i++
		}
	}
	return segments, pre, gap
}

func chooseT(pre []int, tailCandidate int) (int, bool) {
	g := len(pre)
	Tmax := pre[0]
	if tailCandidate < Tmax {
		Tmax = tailCandidate
	}
	for i := 1; i < g; i++ {
		val := pre[i] - 1
		if val < 0 {
			return 0, false
		}
		if val < Tmax {
			Tmax = val
		}
	}
	if Tmax < 0 {
		return 0, false
	}

	blocked := make([]bool, Tmax+1)
	for i := 0; i < g; i++ {
		base := 0
		if i > 0 {
			base = 1
		}
		val := pre[i] - base - 1
		if 0 <= val && val <= Tmax {
			blocked[val] = true
		}
	}
	tailVal := tailCandidate - 1
	if 0 <= tailVal && tailVal <= Tmax {
		blocked[tailVal] = true
	}

	for cand := Tmax; cand >= 0; cand-- {
		if blocked[cand] {
			continue
		}
		if !parityOK(cand, pre, tailCandidate) {
			continue
		}
		ok := true
		for i := 0; i < g; i++ {
			base := 0
			if i > 0 {
				base = 1
			}
			if pre[i] < cand+base {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if tailCandidate < cand {
			continue
		}
		return cand, true
	}
	return 0, false
}

func parityOK(T int, pre []int, tailCandidate int) bool {
	if T != 1 {
		return true
	}
	tailAmt := tailCandidate - 1
	if tailAmt < 0 {
		return false
	}
	if tailAmt > 0 && tailAmt%2 == 1 {
		return false
	}
	for i := range pre {
		base := 0
		if i > 0 {
			base = 1
		}
		x := pre[i] - 1 - base
		if x < 0 {
			return false
		}
		if x > 0 && x%2 == 1 {
			return false
		}
	}
	return true
}

func addZeroBlocks(profile *[]int, amount int, T int) bool {
	if amount == 0 {
		return true
	}
	if amount < 0 || T == 0 {
		return false
	}
	if T == 1 {
		if amount%2 == 1 {
			return false
		}
		for amount > 0 {
			*profile = append(*profile, 1)
			amount -= 2
		}
		return true
	}
	for amount > 0 {
		take := T + 1
		if take > amount {
			take = amount
		}
		if take < 2 {
			return false
		}
		*profile = append(*profile, take-1)
		amount -= take
	}
	return true
}

func buildMask(n int, profile []int) (string, bool) {
	k := len(profile)
	if k == 0 {
		b := make([]byte, n)
		for i := range b {
			b[i] = '_'
		}
		return string(b), true
	}

	earliest := make([]int, k)
	pos := 0
	for i, v := range profile {
		earliest[i] = pos
		pos += v
		if pos > n {
			return "", false
		}
		if i+1 < k {
			pos++
		}
	}
	if pos > n {
		return "", false
	}

	latest := make([]int, k)
	pos = n
	for i := k - 1; i >= 0; i-- {
		pos -= profile[i]
		if pos < 0 {
			return "", false
		}
		latest[i] = pos
		if i > 0 {
			pos--
			if pos < 0 {
				return "", false
			}
		}
	}

	res := make([]byte, n)
	for i := range res {
		res[i] = '_'
	}
	for i := 0; i < k; i++ {
		diff := latest[i] - earliest[i]
		forced := profile[i] - diff
		if forced <= 0 {
			continue
		}
		start := latest[i]
		end := start + forced
		if end > n {
			return "", false
		}
		for j := start; j < end; j++ {
			res[j] = '#'
		}
	}
	return string(res), true
}
