package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func addMod(a, b int64) int64 {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
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

		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		dpHonest := map[int]int64{0: 1} // last person honest, c liars so far
		dpLiar := map[int]int64{}       // last person liar

		for _, val := range a {
			nextHonest := make(map[int]int64)
			nextLiar := make(map[int]int64)

			for c, ways := range dpHonest {
				if val == c {
					nextHonest[c] = addMod(nextHonest[c], ways)
				}
				if c+1 <= n {
					nextLiar[c+1] = addMod(nextLiar[c+1], ways)
				}
			}
			for c, ways := range dpLiar {
				if val == c {
					nextHonest[c] = addMod(nextHonest[c], ways)
				}
			}

			dpHonest = nextHonest
			dpLiar = nextLiar
		}

		var ans int64
		for _, v := range dpHonest {
			ans = addMod(ans, v)
		}
		for _, v := range dpLiar {
			ans = addMod(ans, v)
		}

		fmt.Fprintln(out, ans)
	}
}
