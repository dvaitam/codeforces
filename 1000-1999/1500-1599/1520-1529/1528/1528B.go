package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	d := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j += i {
			d[j]++
		}
	}

	var ans int64 = 1
	for i := 2; i <= n; i++ {
		ans = (2*ans + int64(d[i]-d[i-1])) % mod
	}

	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(out, ans)
}
