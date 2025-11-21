package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func solveCase(a []int64) int64 {
	dpSame := int64(1)
	shift := int64(0)
	// map holds counts for states with exactly two equal values.
	// Keys are stored as (actualDifference XOR shift).
	states := make(map[int64]int64)

	for _, val := range a {
		keyOld := val ^ shift
		cnt := states[keyOld]
		if cnt != 0 {
			delete(states, keyOld)
		}

		dpSameNew := cnt % mod

		shift ^= val

		keyNew := val ^ shift
		add := (3*dpSame + 2*cnt) % mod
		states[keyNew] = (states[keyNew] + add) % mod

		dpSame = dpSameNew
	}

	ans := dpSame % mod
	for _, v := range states {
		ans += v
		if ans >= mod {
			ans -= mod
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, solveCase(a))
	}
}
