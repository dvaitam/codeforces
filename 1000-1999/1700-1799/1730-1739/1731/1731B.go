package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007
const inv6 int64 = 166666668

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		x := n % mod
		t1 := (4 * x % mod) * x % mod
		t2 := (3 * x) % mod
		term := (t1 + t2 - 1) % mod
		if term < 0 {
			term += mod
		}
		ans := x * term % mod
		ans = ans * inv6 % mod
		ans = ans * 2022 % mod
		fmt.Fprintln(out, ans)
	}
}
