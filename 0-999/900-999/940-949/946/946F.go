package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

type trans struct {
	cnt [][]int64
	sum [][]int64
}

func buildNext(s string) [][]int {
	n := len(s)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[j] != s[i] {
			j = pi[j-1]
		}
		if s[j] == s[i] {
			j++
		}
		pi[i] = j
	}
	nxt := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		nxt[i] = make([]int, 2)
		for b := 0; b < 2; b++ {
			k := i
			for k > 0 && (k == n || s[k] != byte('0'+b)) {
				k = pi[k-1]
			}
			if k < n && s[k] == byte('0'+b) {
				k++
			}
			nxt[i][b] = k
		}
	}
	return nxt
}

func charTrans(c byte, nxt [][]int, n int) trans {
	t := trans{cnt: make([][]int64, n+1), sum: make([][]int64, n+1)}
	for i := 0; i <= n; i++ {
		t.cnt[i] = make([]int64, n+1)
		t.sum[i] = make([]int64, n+1)
	}
	idx := int(c - '0')
	for i := 0; i <= n; i++ {
		t.cnt[i][i] = 1
		j2 := nxt[i][idx]
		t.cnt[i][j2] = (t.cnt[i][j2] + 1) % mod
		if j2 == n {
			t.sum[i][j2] = (t.sum[i][j2] + 1) % mod
		}
	}
	return t
}

func compose(a, b trans, n int) trans {
	r := trans{cnt: make([][]int64, n+1), sum: make([][]int64, n+1)}
	for i := 0; i <= n; i++ {
		r.cnt[i] = make([]int64, n+1)
		r.sum[i] = make([]int64, n+1)
	}
	for i := 0; i <= n; i++ {
		for m := 0; m <= n; m++ {
			if a.cnt[i][m] == 0 && a.sum[i][m] == 0 {
				continue
			}
			for k := 0; k <= n; k++ {
				if b.cnt[m][k] == 0 && b.sum[m][k] == 0 {
					continue
				}
				r.cnt[i][k] = (r.cnt[i][k] + a.cnt[i][m]*b.cnt[m][k]) % mod
				val := (a.sum[i][m]*b.cnt[m][k] + a.cnt[i][m]*b.sum[m][k]) % mod
				r.sum[i][k] = (r.sum[i][k] + val) % mod
			}
		}
	}
	return r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, x int
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)
	nxt := buildNext(s)

	transArr := make([]trans, x+1)
	if x >= 0 {
		transArr[0] = charTrans('0', nxt, n)
	}
	if x >= 1 {
		transArr[1] = charTrans('1', nxt, n)
	}
	for i := 2; i <= x; i++ {
		transArr[i] = compose(transArr[i-1], transArr[i-2], n)
	}
	ans := int64(0)
	for k := 0; k <= n; k++ {
		ans = (ans + transArr[x].sum[0][k]) % mod
	}
	fmt.Fprintln(writer, ans)
}
