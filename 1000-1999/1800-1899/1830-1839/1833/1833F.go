package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int = 1000000007

func modPow(a, e int) int {
	res := 1
	x := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(x) % int64(MOD))
		}
		x = int(int64(x) * int64(x) % int64(MOD))
		e >>= 1
	}
	return res
}

func modInv(a int) int {
	return modPow(a, MOD-2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		freq := make(map[int]int)
		for _, v := range arr {
			freq[v]++
		}

		unique := make([]int, 0, len(freq))
		for v := range freq {
			unique = append(unique, v)
		}
		sort.Ints(unique)

		invCache := make(map[int]int)
		getInv := func(x int) int {
			if val, ok := invCache[x]; ok {
				return val
			}
			val := modInv(x)
			invCache[x] = val
			return val
		}

		ans := 0
		prod := 1
		left := 0
		for right := 0; right < len(unique); right++ {
			cntR := freq[unique[right]]
			prod = int(int64(prod) * int64(cntR) % int64(MOD))

			for unique[right]-unique[left] >= m {
				cntL := freq[unique[left]]
				prod = int(int64(prod) * int64(getInv(cntL)) % int64(MOD))
				left++
			}

			for right-left+1 > m {
				cntL := freq[unique[left]]
				prod = int(int64(prod) * int64(getInv(cntL)) % int64(MOD))
				left++
			}

			if right-left+1 == m && unique[right]-unique[left] < m {
				ans += prod
				if ans >= MOD {
					ans -= MOD
				}
			}
		}

		fmt.Fprintln(writer, ans)
	}
}
