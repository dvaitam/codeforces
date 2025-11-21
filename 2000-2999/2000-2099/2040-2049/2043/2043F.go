package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxVal = 50
	mod    = 998244353
	inf    = int(1e9)
)

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
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
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

var fact, invFact []int

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(a) % mod)
		}
		a = int(int64(a) * int64(a) % mod)
		e >>= 1
	}
	return res
}

func initComb(n int) {
	fact = make([]int, n+1)
	invFact = make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = int(int64(invFact[i]) * int64(i) % mod)
	}
}

func comb(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	return int(int64(fact[n]) * int64(invFact[k]) % mod * int64(invFact[n-k]) % mod)
}

func main() {
	in := newScanner()
	n := in.nextInt()
	q := in.nextInt()

	initComb(n)

	pref := make([][maxVal + 1]int, n+1)
	for i := 1; i <= n; i++ {
		x := in.nextInt()
		pref[i] = pref[i-1]
		pref[i][x]++
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	size := [2][64]int{}
	ways := [2][64]int{}

	for ; q > 0; q-- {
		l := in.nextInt()
		r := in.nextInt()
		lenSeg := r - l + 1

		counts := [maxVal + 1]int{}
		for v := 0; v <= maxVal; v++ {
			counts[v] = pref[r][v] - pref[l-1][v]
		}

		for p := 0; p < 2; p++ {
			for m := 0; m < 64; m++ {
				size[p][m] = inf
				ways[p][m] = 0
			}
		}
		size[0][0] = 0
		ways[0][0] = 1

		for val := 0; val <= maxVal; val++ {
			c := counts[val]
			if c == 0 {
				continue
			}
			limit := c
			if limit > 7 {
				limit = 7 // minimal dependent set size in 6-bit space is at most 7
			}

			var nSize [2][64]int
			var nWays [2][64]int
			for p := 0; p < 2; p++ {
				for m := 0; m < 64; m++ {
					nSize[p][m] = inf
					nWays[p][m] = 0
				}
			}

			for picked := 0; picked < 2; picked++ {
				for mask := 0; mask < 64; mask++ {
					if size[picked][mask] == inf {
						continue
					}
					for k := 0; k <= limit; k++ {
						if k > c {
							break
						}
						newPicked := picked
						if k > 0 {
							newPicked = 1
						}
						newMask := mask
						if k&1 == 1 {
							newMask ^= val
						}
						newSize := size[picked][mask] + k
						addWays := int(int64(ways[picked][mask]) * int64(comb(c, k)) % mod)

						if newSize < nSize[newPicked][newMask] {
							nSize[newPicked][newMask] = newSize
							nWays[newPicked][newMask] = addWays
						} else if newSize == nSize[newPicked][newMask] {
							nWays[newPicked][newMask] += addWays
							if nWays[newPicked][newMask] >= mod {
								nWays[newPicked][newMask] -= mod
							}
						}
					}
				}
			}
			size = nSize
			ways = nWays
		}

		ansSize := size[1][0]
		if ansSize == inf {
			fmt.Fprintln(out, -1)
			continue
		}
		removed := lenSeg - ansSize
		fmt.Fprintln(out, removed, ways[1][0]%mod)
	}
}
