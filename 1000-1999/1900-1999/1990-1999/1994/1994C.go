package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)  // 1 MiB input buffer
	out := bufio.NewWriterSize(os.Stdout, 1<<20) // 1 MiB output buffer
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)

		a := make([]int64, n+2) // 1‑indexed; a[0] unused
		b := make([]int64, n+2) // prefix sums, b[0]=0
		f := make([]int64, n+3) // f[n+1] preset to 0

		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
			b[i] = b[i-1] + a[i]
		}

		var ans int64
		j := n
		for i := n; i >= 1; i-- {
			for j > i && b[j-1]-b[i-1] > k {
				j--
			}
			if b[j]-b[i-1] > k { // can’t extend past j
				f[i] = int64(j-i) + f[j+1]
			} else { // whole suffix fits
				f[i] = int64(n - i + 1)
			}
			ans += f[i]
		}
		fmt.Fprintln(out, ans)
	}
}

