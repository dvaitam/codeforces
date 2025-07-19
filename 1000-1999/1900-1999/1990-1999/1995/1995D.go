package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, c, k int
		var s string
		fmt.Fscan(in, &n, &c, &k, &s)

		tag := make([]bool, 1<<uint(c))
		sum := make([][]int, c)
		for j := 0; j < c; j++ {
			sum[j] = make([]int, n+1)
			for i := 1; i <= n; i++ {
				val := 0
				if int(s[i-1]-'A') == j {
					val = 1
				}
				sum[j][i] = sum[j][i-1] + val
			}
		}
		for i := 1; i <= n-k+1; i++ {
			t := 0
			for j := 0; j < c; j++ {
				if sum[j][i+k-1]-sum[j][i-1] == 0 {
					t |= 1 << uint(j)
				}
			}
			tag[t] = true
		}
		for i := (1 << uint(c)) - 1; i > 0; i-- {
			if tag[i] {
				for j := 0; j < c; j++ {
					if (i>>uint(j))&1 != 0 {
						tag[i^(1<<uint(j))] = true
					}
				}
			}
		}
		ans := int(1e9)
		last := int(s[len(s)-1] - 'A')
		for i := 0; i < (1 << uint(c)); i++ {
			if !tag[i] && ((i>>uint(last))&1) != 0 {
				cnt := bits.OnesCount(uint(i))
				ans = min(ans, cnt)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
