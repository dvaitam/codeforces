package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const length = 15

func parseHex(s string) uint64 {
	v, _ := strconv.ParseUint(s, 16, 64)
	return v
}

func countUpto(limit uint64) uint64 {
	var digits [length]int
	for i := length - 1; i >= 0; i-- {
		digits[i] = int(limit & 0xF)
		limit >>= 4
	}
	var total uint64
	for m := 0; m < 16; m++ {
		r := m & 3
		p := m >> 2
		if p >= length {
			continue
		}
		target := length - 1 - p
		var dp [2][2][2]uint64
		dp[0][0][1] = 1
		for i := 0; i < length; i++ {
			var ndp [2][2][2]uint64
			for eq := 0; eq < 2; eq++ {
				for bit := 0; bit < 2; bit++ {
					for tight := 0; tight < 2; tight++ {
						val := dp[eq][bit][tight]
						if val == 0 {
							continue
						}
						upper := m
						if tight == 1 && digits[i] < upper {
							upper = digits[i]
						}
						for d := 0; d <= upper; d++ {
							neEq := eq
							if d == m {
								neEq = 1
							}
							neBit := bit
							if i == target {
								if (d>>r)&1 == 1 {
									neBit = 1
								} else {
									neBit = 0
								}
							}
							neTight := 0
							if tight == 1 && d == digits[i] {
								neTight = 1
							}
							ndp[neEq][neBit][neTight] += val
						}
					}
				}
			}
			dp = ndp
		}
		total += dp[1][1][0] + dp[1][1][1]
	}
	return total
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var lStr, rStr string
		fmt.Fscan(reader, &lStr, &rStr)
		l := parseHex(lStr)
		r := parseHex(rStr)
		ans := countUpto(r)
		if l > 0 {
			ans -= countUpto(l - 1)
		}
		fmt.Fprintln(writer, ans)
	}
}
