package main

import (
	"bufio"
	"fmt"
	"os"
)

func countPairs(S, T, L int) int64 {
	if L < 0 {
		return 0
	}
	leUpper := S
	if leUpper > L {
		leUpper = L
	}
	if leUpper < 0 {
		return 0
	}
	if T >= L {
		cnt := int64(leUpper + 1)
		sumLe := int64(leUpper) * int64(leUpper+1) / 2
		return cnt*int64(L+1) - sumLe
	}
	limit1 := L - T
	r1End := leUpper
	if r1End > limit1 {
		r1End = limit1
	}
	var res int64 = 0
	if r1End >= 0 {
		res += int64(r1End+1) * int64(T+1)
	}
	start2 := limit1 + 1
	if start2 < 0 {
		start2 = 0
	}
	if start2 <= leUpper {
		cnt := int64(leUpper - start2 + 1)
		sumLe := int64(start2+leUpper) * cnt / 2
		res += cnt*int64(L+1) - sumLe
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		pos := make([]int, n)
		for i, v := range p {
			pos[v] = i
		}
		Larr := make([]int, n+1)
		Rarr := make([]int, n+1)
		mn := n
		mx := -1
		for x := 1; x <= n; x++ {
			idx := pos[x-1]
			if idx < mn {
				mn = idx
			}
			if idx > mx {
				mx = idx
			}
			Larr[x] = mn
			Rarr[x] = mx
		}
		var ans int64
		for x := 1; x < n; x++ {
			left := Larr[x]
			right := Rarr[x]
			base := right - left + 1
			if base > 2*x {
				continue
			}
			px := pos[x]
			if px >= left && px <= right {
				continue
			}
			var S, Tspace int
			if px < left {
				S = left - px - 1
				Tspace = n - 1 - right
			} else {
				S = left
				Tspace = px - right - 1
			}
			Lmax := 2*x - base
			ans += countPairs(S, Tspace, Lmax)
		}
		ans++ // whole array for mex = n
		fmt.Fprintln(out, ans)
	}
}
