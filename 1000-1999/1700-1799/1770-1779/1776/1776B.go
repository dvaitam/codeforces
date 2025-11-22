package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var n int
	var h int64
	if _, err := fmt.Fscan(in, &n, &h); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	f := make([][]int64, n)
	for i := range f {
		f[i] = make([]int64, n)
	}

	for l := n - 1; l >= 0; l-- {
		for r := l + 1; r < n; r++ {
			// Split the interval and take the best combination.
			for k := l; k < r; k++ {
				if f[l][k]+f[k+1][r] > f[l][r] {
					f[l][r] = f[l][k] + f[k+1][r]
				}
			}

			// Extend coverage if the two endpoints can be merged.
			if a[l]+h >= a[r]-h {
				diff := a[r] - a[l]
				f[l][r] += h - (diff-1)/2
			}
		}
	}

	ans := int64(n)*h - f[0][n-1]
	fmt.Println(ans)
}
