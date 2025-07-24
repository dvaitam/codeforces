package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 2e18

func mul(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	if a > INF/b {
		return INF
	}
	return a * b
}

func add(a, b int64) int64 {
	if a > INF-b {
		return INF
	}
	return a + b
}

func solveCase(n int, k int64, writer *bufio.Writer) {
	k--
	cycl := make([]int64, n+1)
	cnt := make([]int64, n+1)
	cycl[0], cycl[1] = 1, 1
	for i := 2; i <= n; i++ {
		cycl[i] = mul(cycl[i-1], int64(i-1))
	}
	cnt[n] = 1
	for i := n - 1; i >= 0; i-- {
		var s int64
		for val := i; val < n; val++ {
			length := val - i + 1
			s = add(s, mul(cnt[i+length], cycl[length-1]))
		}
		cnt[i] = s
	}
	if cnt[0] <= k {
		fmt.Fprintln(writer, -1)
		return
	}
	used := make([]bool, n)
	ans := make([]int, n)
	for i := range ans {
		ans[i] = -1
	}
	for i := 0; i < n; i++ {
		for val := i; val < n; val++ {
			length := val - i + 1
			cur := mul(cnt[i+length], cycl[length-1])
			if cur <= k {
				k -= cur
				continue
			}
			ans[i] = val
			used[val] = true
			for j := i + 1; j < i+length; j++ {
				lft := length - (j - i) - 1
				for nval := i; nval <= val; nval++ {
					if used[nval] || nval == j {
						continue
					}
					if j != i+length-1 {
						tmp := ans[nval]
						for tmp != -1 && tmp != j {
							tmp = ans[tmp]
						}
						if tmp == j {
							continue
						}
					}
					cur = mul(cnt[i+length], cycl[lft])
					if cur <= k {
						k -= cur
						continue
					}
					ans[j] = nval
					used[nval] = true
					break
				}
			}
			i += length - 1
			break
		}
	}
	for idx, v := range ans {
		if idx > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v+1)
	}
	fmt.Fprintln(writer)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		solveCase(n, k, out)
	}
}
