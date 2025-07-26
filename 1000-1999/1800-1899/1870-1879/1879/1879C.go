package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	const maxN = 200000
	fact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}

	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		dp0Len, dp1Len := 0, 0
		var dp0Cnt, dp1Cnt int64
		for i := 0; i < n; i++ {
			c := s[i]
			if c == '0' {
				newLen := 1
				var newCnt int64 = 1
				if dp1Len+1 > newLen {
					newLen = dp1Len + 1
					newCnt = dp1Cnt
				} else if dp1Len+1 == newLen {
					newCnt = (newCnt + dp1Cnt) % mod
				}
				if newLen > dp0Len {
					dp0Len = newLen
					dp0Cnt = newCnt % mod
				} else if newLen == dp0Len {
					dp0Cnt = (dp0Cnt + newCnt) % mod
				}
			} else {
				newLen := 1
				var newCnt int64 = 1
				if dp0Len+1 > newLen {
					newLen = dp0Len + 1
					newCnt = dp0Cnt
				} else if dp0Len+1 == newLen {
					newCnt = (newCnt + dp0Cnt) % mod
				}
				if newLen > dp1Len {
					dp1Len = newLen
					dp1Cnt = newCnt % mod
				} else if newLen == dp1Len {
					dp1Cnt = (dp1Cnt + newCnt) % mod
				}
			}
		}
		L := dp0Len
		cnt := dp0Cnt % mod
		if dp1Len > L {
			L = dp1Len
			cnt = dp1Cnt % mod
		} else if dp1Len == L {
			cnt = (cnt + dp1Cnt) % mod
		}
		deletions := n - L
		ans := cnt * fact[deletions] % mod
		fmt.Fprintln(writer, deletions, ans)
	}
}
