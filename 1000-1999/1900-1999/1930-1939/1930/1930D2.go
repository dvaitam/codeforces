package main

import (
	"bufio"
	"fmt"
	"os"
)

func computeSum(s string) int64 {
	var cnt [5]int64
	var sum [5]int64
	var total int64
	for i := 0; i < len(s); i++ {
		ch := s[i]
		var newCnt [5]int64
		var newSum [5]int64
		if ch == '0' {
			newCnt[0] = 1
		} else {
			newCnt[1] = 1
			newSum[1] = 1
		}
		for st := 0; st < 5; st++ {
			c := cnt[st]
			if c == 0 {
				continue
			}
			if ch == '0' {
				switch st {
				case 0:
					newCnt[0] += c
					newSum[0] += sum[st]
				case 1:
					newCnt[2] += c
					newSum[2] += sum[st]
				case 2:
					newCnt[3] += c
					newSum[3] += sum[st]
				default: // 3 or 4
					newCnt[4] += c
					newSum[4] += sum[st]
				}
			} else { // ch == '1'
				switch st {
				case 0:
					newCnt[1] += c
					newSum[1] += sum[st] + c
				case 1:
					newCnt[2] += c
					newSum[2] += sum[st]
				case 2:
					newCnt[3] += c
					newSum[3] += sum[st]
				case 3, 4:
					newCnt[1] += c
					newSum[1] += sum[st] + c
				}
			}
		}
		cnt = newCnt
		sum = newSum
		for j := 0; j < 5; j++ {
			total += sum[j]
		}
	}
	return total
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, computeSum(s))
	}
}
