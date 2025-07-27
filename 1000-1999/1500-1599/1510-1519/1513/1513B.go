package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN int = 200000

var fact [MAXN + 1]int64

func init() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		andAll := arr[0]
		for i := 1; i < n; i++ {
			andAll &= arr[i]
		}
		cnt := 0
		for i := 0; i < n; i++ {
			if arr[i] == andAll {
				cnt++
			}
		}
		if cnt < 2 {
			fmt.Fprintln(writer, 0)
			continue
		}
		ans := int64(cnt) * int64(cnt-1) % MOD
		ans = ans * fact[n-2] % MOD
		fmt.Fprintln(writer, ans)
	}
}
