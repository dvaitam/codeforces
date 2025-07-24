package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var mod int64
	if _, err := fmt.Fscan(reader, &n, &mod); err != nil {
		return
	}

	fac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}

	var ans int64
	for k := 1; k <= n; k++ {
		term := fac[k] * fac[n-k] % mod
		length := int64(n - k + 1)
		term = term * length % mod
		term = term * length % mod
		ans = (ans + term) % mod
	}

	fmt.Fprintln(writer, ans)
}
