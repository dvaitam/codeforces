package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	p := make([][]float64, n+1)
	for i := 1; i <= n; i++ {
		p[i] = make([]float64, n+1)
		for j := 1; j <= n; j++ {
			fmt.Fscan(reader, &p[i][j])
		}
	}

	maxMask := (1 << uint(n)) - 1
	dp := make([][]float64, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]float64, maxMask+1)
	}

	for mask := 0; mask <= maxMask; mask++ {
		bitsCnt := bits.OnesCount(uint(mask))
		for c := 0; c <= n; c++ {
			if c != 0 && (mask&(1<<(c-1))) != 0 {
				dp[c][mask] = 0
				continue
			}
			if bitsCnt == 0 {
				if c == 1 {
					dp[c][mask] = 1.0
				} else {
					dp[c][mask] = 0.0
				}
				continue
			}

			if c != 0 {
				mx := 0.0
				for ib := 0; ib < n; ib++ {
					if mask&(1<<ib) == 0 {
						continue
					}
					x := ib + 1
					newMask := mask &^ (1 << ib)
					prob := p[c][x]*dp[c][newMask] + p[x][c]*dp[x][newMask]
					if prob > mx {
						mx = prob
					}
				}
				dp[c][mask] = mx
			} else {
				if bitsCnt == 1 {
					var x int
					for ib := 0; ib < n; ib++ {
						if mask&(1<<ib) != 0 {
							x = ib + 1
							break
						}
					}
					if x == 1 {
						dp[0][mask] = 1.0
					} else {
						dp[0][mask] = 0.0
					}
				} else if bitsCnt >= 2 {
					mx := 0.0
					for ib := 0; ib < n; ib++ {
						if mask&(1<<ib) == 0 {
							continue
						}
						for jb := ib + 1; jb < n; jb++ {
							if mask&(1<<jb) == 0 {
								continue
							}
							x := ib + 1
							y := jb + 1
							newMask := mask &^ ((1 << ib) | (1 << jb))
							prob := p[x][y]*dp[x][newMask] + p[y][x]*dp[y][newMask]
							if prob > mx {
								mx = prob
							}
						}
					}
					dp[0][mask] = mx
				} else {
					dp[0][mask] = 0.0
				}
			}
		}
	}

	fmt.Printf("%.10f\n", dp[0][maxMask])
}
