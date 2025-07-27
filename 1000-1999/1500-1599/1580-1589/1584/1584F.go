package main

import (
	"bufio"
	"fmt"
	"os"
)

const Letters = 52
const MaxMask = 1 << 10

var n int
var pos [10][Letters][2]int
var dp [Letters + 1][MaxMask]int
var vis [Letters + 1][MaxMask]bool
var nextC [Letters + 1][MaxMask]int
var nextM [Letters + 1][MaxMask]int

func idx(c byte) int {
	if c >= 'a' && c <= 'z' {
		return int(c - 'a')
	}
	return 26 + int(c-'A')
}

func toChar(i int) byte {
	if i < 26 {
		return byte('a' + i)
	}
	return byte('A' + i - 26)
}

func nextState(curC, mask, d int) (int, bool) {
	nextMask := 0
	for i := 0; i < n; i++ {
		prev := -1
		if curC != Letters {
			prev = pos[i][curC][(mask>>i)&1]
			if prev == -1 {
				return 0, false
			}
		}
		p1 := pos[i][d][0]
		p2 := pos[i][d][1]
		best := int(1e9)
		bit := 0
		if p1 != -1 && p1 > prev && p1 < best {
			best = p1
			bit = 0
		}
		if p2 != -1 && p2 > prev && p2 < best {
			best = p2
			bit = 1
		}
		if best == int(1e9) {
			return 0, false
		}
		if bit == 1 {
			nextMask |= 1 << i
		}
	}
	return nextMask, true
}

func dfs(c, mask int) int {
	if vis[c][mask] {
		return dp[c][mask]
	}
	vis[c][mask] = true
	best := 0
	bC, bM := -1, 0
	for d := 0; d < Letters; d++ {
		nm, ok := nextState(c, mask, d)
		if !ok {
			continue
		}
		length := dfs(d, nm)
		if length > best {
			best = length
			bC = d
			bM = nm
		}
	}
	dp[c][mask] = best + 1
	nextC[c][mask] = bC
	nextM[c][mask] = bM
	return dp[c][mask]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		fmt.Fscan(in, &n)
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &strs[i])
		}
		for i := 0; i < n; i++ {
			for j := 0; j < Letters; j++ {
				pos[i][j][0], pos[i][j][1] = -1, -1
			}
			for k := 0; k < len(strs[i]); k++ {
				id := idx(strs[i][k])
				if pos[i][id][0] == -1 {
					pos[i][id][0] = k
				} else {
					pos[i][id][1] = k
				}
			}
		}
		for i := 0; i <= Letters; i++ {
			for m := 0; m < (1 << n); m++ {
				dp[i][m] = 0
				vis[i][m] = false
				nextC[i][m] = -1
				nextM[i][m] = 0
			}
		}
		maxLen := dfs(Letters, 0) - 1
		fmt.Println(maxLen)
		res := make([]byte, 0, maxLen)
		c, m := Letters, 0
		for {
			nc := nextC[c][m]
			if nc == -1 {
				break
			}
			nm := nextM[c][m]
			res = append(res, toChar(nc))
			c, m = nc, nm
		}
		fmt.Println(string(res))
	}
}
