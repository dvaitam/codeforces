package main

import (
	"bufio"
	"fmt"
	"os"
)

const base uint64 = 911382323

func getHash(h, pow []uint64, l, r int) uint64 {
	return h[r] - h[l-1]*pow[r-l+1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		var s string
		fmt.Fscan(in, &s)

		ps1 := make([]int, n+1)
		for i := 2; i <= n; i++ {
			ps1[i] = ps1[i-1]
			if s[i-1] != s[i-2] {
				ps1[i]++
			}
		}

		ps2 := make([]int, n+1)
		for i := 3; i <= n; i++ {
			ps2[i] = ps2[i-1]
			if s[i-1] != s[i-3] {
				ps2[i]++
			}
		}

		pow := make([]uint64, n+1)
		hf := make([]uint64, n+1)
		hr := make([]uint64, n+1)
		pow[0] = 1
		for i := 1; i <= n; i++ {
			pow[i] = pow[i-1] * base
			ch := uint64(s[i-1] - 'a' + 1)
			hf[i] = hf[i-1]*base + ch
			chR := uint64(s[n-i] - 'a' + 1)
			hr[i] = hr[i-1]*base + chR
		}

		isPal := func(l, r int) bool {
			hfwd := getHash(hf, pow, l, r)
			rl := n - r + 1
			rr := n - l + 1
			hrev := getHash(hr, pow, rl, rr)
			return hfwd == hrev
		}

		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			m := r - l + 1
			total := m * (m + 1) / 2
			if ps1[r]-ps1[l] == 0 {
				fmt.Fprintln(out, 0)
				continue
			}
			if s[l-1] != s[l] && ps2[r]-ps2[l+1] == 0 {
				L := (m + 1) / 2
				ans := total - L*L
				fmt.Fprintln(out, ans)
				continue
			}
			if isPal(l, r) {
				ans := total - (m + 1)
				fmt.Fprintln(out, ans)
				continue
			}
			fmt.Fprintln(out, total-1)
		}
	}
}
