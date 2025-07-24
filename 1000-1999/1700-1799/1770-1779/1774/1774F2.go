package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const INF int64 = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	opsType := make([]int, n)
	opsVal := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &opsType[i])
		if opsType[i] != 3 {
			fmt.Fscan(in, &opsVal[i])
		}
	}

	prefDamage := make([]int64, n+1)
	for i := 0; i < n; i++ {
		d := prefDamage[i]
		switch opsType[i] {
		case 1:
			prefDamage[i+1] = d
		case 2:
			d += opsVal[i]
			if d > INF {
				d = INF
			}
			prefDamage[i+1] = d
		case 3:
			d *= 2
			if d > INF {
				d = INF
			}
			prefDamage[i+1] = d
		}
	}
	totalDamage := prefDamage[n]

	total3 := 0
	cnt3Prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		cnt3Prefix[i+1] = cnt3Prefix[i]
		if opsType[i] == 3 {
			cnt3Prefix[i+1]++
			total3++
		}
	}

	threeVal := make([]int64, total3)
	idx := 0
	for i := 0; i < n; i++ {
		if opsType[i] == 3 {
			threeVal[idx] = prefDamage[i]
			idx++
		}
	}

	pow2 := make([]int64, total3+1)
	pow2[0] = 1
	for i := 1; i <= total3; i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}

	ans := int64(0)
	for i := 0; i < n; i++ {
		if opsType[i] != 1 {
			continue
		}
		x := opsVal[i]
		damageBefore := prefDamage[i]
		cntAfter := total3 - cnt3Prefix[i+1]
		r := totalDamage - damageBefore - x
		if r < 0 {
			ans = (ans + pow2[cntAfter]) % MOD
			continue
		}
		indexMod := int64(0)
		for j := total3 - 1; j >= cnt3Prefix[i+1]; j-- {
			w := threeVal[j]
			if r >= w {
				r -= w
				indexMod = (indexMod*2 + 1) % MOD
			} else {
				indexMod = (indexMod * 2) % MOD
			}
		}
		add := (pow2[cntAfter] - indexMod - 1) % MOD
		if add < 0 {
			add += MOD
		}
		ans = (ans + add) % MOD
	}

	fmt.Println(ans % MOD)
}
