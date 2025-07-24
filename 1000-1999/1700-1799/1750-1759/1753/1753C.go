package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353
const maxn = 200000

var inv [maxn + 1]int64
var pref [maxn + 1]int64

func init() {
	inv[1] = 1
	for i := 2; i <= maxn; i++ {
		inv[i] = mod - mod/int64(i)*inv[mod%int64(i)]%mod
	}
	for i := 1; i <= maxn; i++ {
		pref[i] = (pref[i-1] + inv[i]*inv[i]) % mod
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		ones := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] == 1 {
				ones++
			}
		}
		zeros := n - ones
		k := 0
		for i := 0; i < zeros; i++ {
			if a[i] == 1 {
				k++
			}
		}
		C := int64(n) * int64(n-1) / 2 % mod
		ans := C * pref[k] % mod
		fmt.Fprintln(writer, ans)
	}
}
