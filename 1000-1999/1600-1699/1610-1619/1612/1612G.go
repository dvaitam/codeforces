package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007
const INV2 int64 = (MOD + 1) / 2

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	c := make([]int, m)
	maxC := 0
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &c[i])
		if c[i] > maxC {
			maxC = c[i]
		}
	}

	arrOdd := make([]int64, maxC)
	arrEven := make([]int64, maxC)
	for _, v := range c {
		if v%2 == 1 {
			arrOdd[v-1]++
		} else {
			arrEven[v-1]++
		}
	}

	prefOdd := make([]int64, maxC+1)
	prefEven := make([]int64, maxC+1)
	for t := maxC - 1; t >= 0; t-- {
		prefOdd[t] = prefOdd[t+1] + arrOdd[t]
		prefEven[t] = prefEven[t+1] + arrEven[t]
	}

	fact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}

	pos := int64(1)
	ans := int64(0)
	ways := int64(1)
	for x := -maxC + 1; x <= maxC-1; x++ {
		t := x
		if t < 0 {
			t = -t
		}
		if t >= maxC {
			continue
		}
		var cnt int64
		if t%2 == 0 {
			cnt = prefOdd[t]
		} else {
			cnt = prefEven[t]
		}
		if cnt == 0 {
			continue
		}
		cntMod := cnt % MOD
		posMod := pos % MOD
		sumPosMod := (cntMod * posMod) % MOD
		sumPosMod = (sumPosMod + cntMod*((cnt-1)%MOD)%MOD*INV2%MOD) % MOD
		xMod := int64(x % int(MOD))
		if xMod < 0 {
			xMod += MOD
		}
		ans = (ans + xMod*sumPosMod%MOD) % MOD
		ways = ways * fact[int(cnt)] % MOD
		pos += cnt
	}

	fmt.Fprintf(writer, "%d %d\n", ans, ways)
}
