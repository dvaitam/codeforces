package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func thueMorse(x uint64) uint64 {
	return uint64(bits.OnesCount64(x) & 1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m uint64
		fmt.Fscan(in, &n, &m)
		var ans uint64
		for i := uint64(0); i < m; i++ {
			if thueMorse(i)^thueMorse(i+n) == 1 {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
