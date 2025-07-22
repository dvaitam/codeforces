package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const base uint64 = 911382323

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	pow := make([]uint64, n+1)
	pow[0] = 1
	for i := 1; i <= n; i++ {
		pow[i] = pow[i-1] * base
	}

	pref := make([][]uint64, 26)
	for c := 0; c < 26; c++ {
		pref[c] = make([]uint64, n+1)
		cur := uint64(0)
		for i := 0; i < n; i++ {
			cur = cur * base
			if int(s[i]-'a') == c {
				cur++
			}
			pref[c][i+1] = cur
		}
	}

	for ; m > 0; m-- {
		var x, y, l int
		fmt.Fscan(reader, &x, &y, &l)
		x--
		y--
		a := make([]uint64, 26)
		b := make([]uint64, 26)
		for c := 0; c < 26; c++ {
			pa := pref[c]
			a[c] = pa[x+l] - pa[x]*pow[l]
			b[c] = pa[y+l] - pa[y]*pow[l]
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		ok := true
		for i := 0; i < 26 && ok; i++ {
			if a[i] != b[i] {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
