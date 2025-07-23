package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	aPrev := int64(1)
	dPrev := int64(1)
	dPrevPrev := int64(0)

	for i := 1; i <= n; i++ {
		ai := (2*aPrev - dPrevPrev) % mod
		if ai < 0 {
			ai += mod
		}
		ai = ai * int64(m) % mod
		di := (ai + int64(m-1)*dPrev%mod) % mod
		dPrevPrev, dPrev, aPrev = dPrev, di, ai
	}

	fmt.Fprintln(out, aPrev%mod)
}
