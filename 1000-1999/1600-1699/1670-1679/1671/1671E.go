package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const mod int64 = 998244353

var (
	n       int
	s       string
	letters []byte
	N       int
)

func dfs(idx int) ([]byte, int64) {
	if idx*2 > N {
		return []byte{letters[idx]}, 1
	}
	ls, lc := dfs(idx * 2)
	rs, rc := dfs(idx*2 + 1)
	if bytes.Compare(ls, rs) > 0 {
		ls, rs = rs, ls
		lc, rc = rc, lc
	}
	count := (lc * rc) % mod
	if !bytes.Equal(ls, rs) {
		count = (count * 2) % mod
	}
	res := make([]byte, 1+len(ls)+len(rs))
	res[0] = letters[idx]
	copy(res[1:], ls)
	copy(res[1+len(ls):], rs)
	return res, count
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	fmt.Fscan(in, &s)
	N = (1 << n) - 1
	letters = make([]byte, N+1)
	for i := 1; i <= N; i++ {
		letters[i] = s[i-1]
	}
	_, ans := dfs(1)
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans%mod)
	out.Flush()
}
