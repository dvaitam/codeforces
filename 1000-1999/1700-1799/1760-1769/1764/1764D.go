package main

import (
	"bufio"
	"fmt"
	"os"
)

func modFact(n int, mod int64) int64 {
	res := int64(1)
	for i := 2; i <= n; i++ {
		res = res * int64(i) % mod
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var p int64
	if _, err := fmt.Fscan(in, &n, &p); err != nil {
		return
	}

	factN := modFact(n, p)
	factA := modFact(n/2, p)
	factB := modFact((n+1)/2, p)
	ans := (factN - 2*(factA*factB%p)) % p
	if ans < 0 {
		ans += p
	}
	fmt.Fprintln(out, ans)
}
