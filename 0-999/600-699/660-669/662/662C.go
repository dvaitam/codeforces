package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func fwt(a []int64, invert bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				x := a[i+j]
				y := a[i+j+step]
				a[i+j] = x + y
				a[i+j+step] = x - y
			}
		}
	}
	if invert {
		for i := 0; i < n; i++ {
			a[i] /= int64(n)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &rows[i])
	}

	size := 1 << n
	cnt := make([]int64, size)

	for col := 0; col < m; col++ {
		mask := 0
		for i := 0; i < n; i++ {
			if rows[i][col] == '1' {
				mask |= 1 << i
			}
		}
		cnt[mask]++
	}

	weight := make([]int64, size)
	for mask := 0; mask < size; mask++ {
		k := bits.OnesCount(uint(mask))
		if k > n-k {
			k = n - k
		}
		weight[mask] = int64(k)
	}

	fwt(cnt, false)
	fwt(weight, false)
	for i := 0; i < size; i++ {
		cnt[i] *= weight[i]
	}
	fwt(cnt, true)

	minVal := cnt[0]
	for _, v := range cnt {
		if v < minVal {
			minVal = v
		}
	}
	fmt.Println(minVal)
}
