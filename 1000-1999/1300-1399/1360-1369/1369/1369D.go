package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN = 2000000

var ans [MAXN + 1]int64

func init() {
	ans[1] = 0
	ans[2] = 0
	ans[3] = 4
	ans[4] = 4
	for i := 5; i <= MAXN; i++ {
		ans[i] = (2 * ans[i-1]) % MOD
		switch i % 6 {
		case 3, 5:
			ans[i] = (ans[i] + 4) % MOD
		case 4:
			ans[i] = (ans[i] + MOD - 4) % MOD
		}
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
		fmt.Fprintln(writer, ans[n])
	}
}
