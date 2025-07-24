package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit int64 = 1e18

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int64
	fmt.Fscan(reader, &n)
	var k string
	fmt.Fscan(reader, &k)

	s := k
	L := len(s)
	pow := make([]int64, L+1)
	pow[0] = 1
	for i := 1; i <= L; i++ {
		if pow[i-1] > limit/n {
			pow[i] = limit + 1
		} else {
			pow[i] = pow[i-1] * n
			if pow[i] > limit {
				pow[i] = limit + 1
			}
		}
	}

	dp := make([]map[int]int64, L+1)
	dp[L] = map[int]int64{0: 0}

	for pos := L - 1; pos >= 0; pos-- {
		dp[pos] = make(map[int]int64)
		var val int64
		for length := 1; length <= 9 && pos+length <= L; length++ {
			if length > 1 && s[pos] == '0' {
				break
			}
			digit := int64(s[pos+length-1] - '0')
			val = val*10 + digit
			if val >= n {
				break
			}
			for d, nextVal := range dp[pos+length] {
				if pow[d] > limit/val {
					continue
				}
				candidate := val*pow[d] + nextVal
				if candidate > limit {
					continue
				}
				if cur, ok := dp[pos][d+1]; !ok || candidate < cur {
					dp[pos][d+1] = candidate
				}
			}
		}
	}

	ans := int64(-1)
	for _, v := range dp[0] {
		if ans == -1 || v < ans {
			ans = v
		}
	}
	fmt.Fprintln(writer, ans)
}
