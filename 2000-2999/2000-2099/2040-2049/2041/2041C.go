package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b > '~') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if b == '-' {
		sign = -1
		b, err = fs.r.ReadByte()
	}
	for err == nil && b >= '0' && b <= '9' {
		val = val*10 + int(b-'0')
		b, err = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	n := fs.nextInt()
	cost := make([][][]int64, n)
	for x := 0; x < n; x++ {
		cost[x] = make([][]int64, n)
		for y := 0; y < n; y++ {
			cost[x][y] = make([]int64, n)
			for z := 0; z < n; z++ {
				cost[x][y][z] = int64(fs.nextInt())
			}
		}
	}

	if n == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	nbits := uint(n)
	maskAll := (1 << n) - 1
	dp := map[uint32]int64{0: 0}

	for x := 0; x < n; x++ {
		nextDP := make(map[uint32]int64)
		for key, cur := range dp {
			maskZ := int(key & uint32(maskAll))
			maskY := int(key >> nbits)
			for y := 0; y < n; y++ {
				if (maskY>>y)&1 == 1 {
					continue
				}
				newMaskY := maskY | (1 << y)
				for z := 0; z < n; z++ {
					if (maskZ>>z)&1 == 1 {
						continue
					}
					newMaskZ := maskZ | (1 << z)
					newKey := uint32(newMaskY<<n | newMaskZ)
					val := cur + cost[x][y][z]
					if prev, ok := nextDP[newKey]; !ok || val < prev {
						nextDP[newKey] = val
					}
				}
			}
		}
		dp = nextDP
	}

	finalKey := uint32(maskAll<<n | maskAll)
	if ans, ok := dp[finalKey]; ok {
		fmt.Fprintln(out, ans)
	} else {
		fmt.Fprintln(out, 0)
	}
}
