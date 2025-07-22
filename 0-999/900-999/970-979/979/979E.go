package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, p int
	if _, err := fmt.Fscan(in, &n, &p); err != nil {
		return
	}
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &colors[i])
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}

	dp := map[int]int64{0: 1} // key encodes e0<<2 | e1<<1 | parity
	for idx := n - 1; idx >= 0; idx-- {
		edgesRem := n - idx - 1
		newDP := make(map[int]int64)
		colOptions := []int{colors[idx]}
		if colors[idx] == -1 {
			colOptions = []int{0, 1}
		}
		for key, ways := range dp {
			e0 := (key >> 2) & 1
			e1 := (key >> 1) & 1
			parity := key & 1
			for _, col := range colOptions {
				diffExist := e1
				if col == 1 {
					diffExist = e0
				}
				if diffExist == 1 {
					edgeWays := int64(1)
					if edgesRem > 0 {
						edgeWays = pow2[edgesRem-1]
					}
					for fi := 0; fi <= 1; fi++ {
						ne0, ne1 := e0, e1
						if fi == 1 {
							if col == 0 {
								ne0 = 1
							} else {
								ne1 = 1
							}
						}
						np := parity ^ fi
						nk := (ne0 << 2) | (ne1 << 1) | np
						newDP[nk] = (newDP[nk] + ways*edgeWays) % MOD
					}
				} else {
					edgeWays := pow2[edgesRem]
					ne0, ne1 := e0, e1
					if col == 0 {
						ne0 = 1
					} else {
						ne1 = 1
					}
					np := parity ^ 1
					nk := (ne0 << 2) | (ne1 << 1) | np
					newDP[nk] = (newDP[nk] + ways*edgeWays) % MOD
				}
			}
		}
		dp = newDP
	}

	var result int64
	for key, ways := range dp {
		if key&1 == p {
			result = (result + ways) % MOD
		}
	}
	fmt.Println(result)
}
