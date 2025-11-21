package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt64() int64 {
	var sign int64 = 1
	var val int64
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := newScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(in.nextInt64())
	for ; t > 0; t-- {
		n := int(in.nextInt64())
		_ = in.nextInt64() // x is always 1 in this version

		a := make([]int64, n+2)
		for i := 1; i <= n; i++ {
			a[i] = in.nextInt64()
		}

		pref := make([]int64, n+2)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + a[i]
		}
		// Sparse tables for max(a[i]+pref[i]) and min(pref[i]-a[i+1])
		maxK := 1
		for (1 << maxK) <= n {
			maxK++
		}
		f := make([][]int64, maxK)
		g := make([][]int64, maxK)
		for k := 0; k < maxK; k++ {
			f[k] = make([]int64, n+2)
			g[k] = make([]int64, n+2)
		}
		for i := 1; i <= n; i++ {
			f[0][i] = a[i] + pref[i]
			g[0][i] = pref[i] - a[i+1]
		}
		for k := 1; k < maxK; k++ {
			span := 1 << k
			half := span >> 1
			for i := 1; i+span-1 <= n; i++ {
				f[k][i] = max(f[k-1][i], f[k-1][i+half])
				g[k][i] = min(g[k-1][i], g[k-1][i+half])
			}
		}

		logv := make([]int, n+2)
		for i := 2; i <= n+1; i++ {
			logv[i] = logv[i>>1] + 1
		}

		getMax := func(l, r int) int64 {
			if l > r {
				return -1 << 60
			}
			k := logv[r-l+1]
			return max(f[k][l], f[k][r-(1<<k)+1])
		}
		getMin := func(l, r int) int64 {
			if l > r {
				return 1 << 60
			}
			k := logv[r-l+1]
			return min(g[k][l], g[k][r-(1<<k)+1])
		}

		const inf = int(1e9)
		minr := make([]int, n+1)
		maxr := make([]int, n+1)
		for i := 1; i <= n; i++ {
			minr[i] = inf
			maxr[i] = inf
		}

		for i := 1; i <= n; i++ {
			l, r := i, i
			for {
				flag := false
				L, R := 1, l
				for L < R {
					mid := (L + R) >> 1
					if pref[r] < getMax(mid, R-1) {
						L = mid + 1
					} else {
						R = mid
					}
				}
				if l != L {
					l = L
					flag = true
				}
				if l == 1 || r == n {
					break
				}
				L, R = r, n
				for L < R {
					mid := (L + R) >> 1
					if getMin(L, mid) < pref[l-1] {
						R = mid
					} else {
						L = mid + 1
					}
				}
				if r != L {
					flag = true
					x := L
					L, R = r+1, x
					for L < R {
						mid := (L + R) >> 1
						if pref[mid]-pref[l-1] >= a[l-1] {
							R = mid
						} else {
							L = mid + 1
						}
					}
					r = R
				}
				if !flag {
					break
				}
			}
			if l == 1 {
				minr[i] = r
			}

			l, r = i, i
			for {
				flag := false
				L, R := r, n
				for L < R {
					mid := (L + R) >> 1
					if getMin(L, mid) < pref[l-1] {
						R = mid
					} else {
						L = mid + 1
					}
				}
				if r != L {
					flag = true
					r = R
				}
				L, R = 1, l
				for L < R {
					mid := (L + R) >> 1
					if pref[r] < getMax(mid, R-1) {
						L = mid + 1
					} else {
						R = mid
					}
				}
				if l != L {
					flag = true
					l = L
				}
				if !flag {
					break
				}
			}
			if l == 1 {
				maxr[i] = r
			}
		}

		ans := make([]int, n+2)
		for i := 1; i <= n; i++ {
			if maxr[i] != inf && minr[i] <= maxr[i] {
				ans[minr[i]]++
				ans[maxr[i]+1]--
			}
		}
		for i := 1; i <= n; i++ {
			ans[i] += ans[i-1]
			fmt.Fprint(out, ans[i])
			if i == n {
				fmt.Fprintln(out)
			} else {
				fmt.Fprint(out, " ")
			}
		}
	}
}
