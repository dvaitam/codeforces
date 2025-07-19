package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000009

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	fmt.Fscan(in, &n, &m)
	base := make([]int64, 26)
	base[0] = 233333
	for i := 1; i < 26; i++ {
		base[i] = base[i-1] * base[0] % mod
	}
	H := make([][][]int64, n)
	for i := 0; i < n; i++ {
		H[i] = make([][]int64, m)
		for L := 0; L < m; L++ {
			H[i][L] = make([]int64, m)
		}
	}
	idx := int64(mod)
	cnt := make([]bool, 26)
	ch := make([]int, m)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		for j := 0; j < m; j++ {
			ch[j] = int(s[j] - 'a')
		}
		for L := 0; L < m; L++ {
			for k := 0; k < 26; k++ {
				cnt[k] = false
			}
			odd := 0
			cur := int64(0)
			for R := L; R < m; R++ {
				t := ch[R]
				cnt[t] = !cnt[t]
				if cnt[t] {
					odd++
				} else {
					odd--
				}
				cur += base[t]
				if cur >= mod {
					cur -= mod
				}
				if odd <= 1 {
					H[i][L][R] = cur
				} else {
					H[i][L][R] = idx
					idx++
				}
			}
		}
	}
	sentinel := idx
	idx++
	size := 2*n + 3
	S := make([]int64, size)
	P := make([]int, size)
	S[0] = sentinel
	S[2*n+1] = -1
	S[2*n+2] = -1
	var ans uint64
	for L := 0; L < m; L++ {
		for R := L; R < m; R++ {
			j := 0
			for i := 0; i < n; i++ {
				j++
				S[j] = -1
				j++
				S[j] = H[i][L][R]
			}
			for i := range P {
				P[i] = 0
			}
			C, M := 0, 0
			for i := 2; i <= 2*n; i++ {
				if S[i] < mod {
					if i < M {
						mirror := 2*C - i
						p := P[mirror]
						if p > M-i {
							p = M - i
						}
						P[i] = p
					}
					for i-P[i] >= 0 && i+P[i] < size && S[i-P[i]] == S[i+P[i]] {
						P[i]++
					}
					if i+P[i] > M {
						C = i
						M = i + P[i]
					}
					ans += uint64(P[i] >> 1)
				}
			}
		}
	}
	fmt.Fprintln(out, ans)
}
